package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gagangoku/easyetl/parser"
	"github.com/samber/lo"
)

/*
1. Parse command into a datastructure -
2. Parse all commands and surface errors -
3. Generate execution plan -
4. Execute - local vs clickhouse vs BQ

Observations:
- Transformations that are independent can be executed in separate go routines
- Single column computations can be parallelized in separate go routines
- Columns not needed in the final output need not be stored, they could just be computed on the fly
- only LOOKUP requires column index, which can be deferred to the time of evaluation

Time consuming steps:
- Filter conditions seem to be taking longest right now
- Knowing the types of columns might improve performance (less type conversions)

Materialization:
- When new data flows in, not all previously computed cells need to be recomputed

Sequencing:
- Index only columns that need to be looked up ✅
- Parse into infix first, and do simple DFS solve ✅
- Parallelize Filter command execution into multiple goroutines ✅
- Save state in DB / Disk ✅
- Cartesian to allow multiple columns from same table ✅
- Parallelize single column execution into multiple goroutines
- Topological sorting of all columns and parallelize commands at the same level
- Postfix operations, remove nesting / recursion
- Remove type conversions by strongly typing columns


Checkpointing & Persistence:
Many times you need to debug one particular command. The underlying data hasn't really changed.
Imagine if the dependent columns can be fetched from precomputed store instead of recomputing all the commands from scratch.

ColMetadata:
- formula
- for each table + db terms - last computed time + last checksum
- checksum of own rows
- last computed time

How to check if a column needs to be recomputed:
- the formula changed
- dependencies changed (checksum is different from what we know)

Where to store ?
- file store (local, gcp)
- DB as blob (Redis)
- format (arrow, parquet, own format, protobuf)
- ANY type or use native type ?

*/

type Node struct {
	Ctx       string          `json:"ctx"`
	Type      string          `json:"type,omitempty"`
	Text      string          `json:"text,omitempty"`
	Terminal  ResultValueType `json:"terminal,omitempty"`
	Function  string          `json:"fn,omitempty"`
	Children  []*Node         `json:"children,omitempty"`
	TblDbName string          `json:"tbldb,omitempty"`
	Col       string          `json:"col,omitempty"`
}

/*
Assumptions:
- Array based sequential lookups will be faster than random access. So my execution flow should incorporate postfix

a + b + c => [a b + c +]
MIN(a, b) + MAX(c, d, z) * MIN(MAX(e, f), g) => [a, b, MIN, c, d, z, MAX(3), e, f, MAX, g, MIN, *, +]
LAST ( SPLIT ( a, b ) ) => [a, b, SPLIT, LAST]
*/
type PostfixNode struct {
	Function *string `json:"fn"`
	NumArgs  int     `json:"numArgs"`
	Constant *any    `json:"constant"`
}
type PostfixPlan []PostfixNode

type ColMetdata struct {
	Column         string `json:"column"`
	Formula        string `json:"formula"`
	Deps           []_Dep `json:"deps"`
	Checksum       string `json:"checksum"`
	LastModifiedMs int64  `json:"lastModifiedMs"`
}
type _Dep struct {
	Term           string `json:"term"`
	Checksum       string `json:"checksum"`
	LastModifiedMs int64  `json:"lastModifiedMs"`
}

type Engine struct {
	app          *EtlApp
	dbs          map[string]Table
	tables       map[string]Table
	cmdIdx       int
	cmd          string
	indexes      map[string]ColIndex
	loader       LoaderInf
	metadata     map[string]ColMetdata
	evaluatorCtx *Evaluator

	lhs                 *Node
	lhsIndex            int
	nRowsPerLookupBatch int
	fCtx                *FilterContext
}

func (engine *Engine) Init(app *EtlApp) {
	engine.app = app
	engine.dbs = make(map[string]Table)
	engine.tables = make(map[string]Table)
	engine.cmdIdx = 0
	engine.cmd = ""
	engine.indexes = make(map[string]ColIndex)
	engine.metadata = make(map[string]ColMetdata)
	engine.evaluatorCtx = NewEvaluator(engine)
	engine.nRowsPerLookupBatch = 10000

	fl := &FileLoader{}
	fl.Init(app.backupDir, app.appId)
	engine.loader = fl
}

func (engine *Engine) isColMetadataValid(metadata ColMetdata, formula string) bool {
	if metadata.Formula != formula || metadata.Checksum == "" {
		return false
	}
	for _, dep := range metadata.Deps {
		// NOTE: Not checking last modified time here since checksum is good enough
		if c, found := engine.metadata[dep.Term]; !found || c.Checksum != dep.Checksum {
			return false
		}
	}
	return true
}

func (engine *Engine) Parse(cmd string) (*Node, *Node, error) {
	is := antlr.NewInputStream(cmd)

	// Create the Lexer
	lexer := parser.NewEllieLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewEllieParser(stream)

	// Parse the expression
	res, err := engine.dfsParse(p.Prog())
	if err != nil {
		return nil, nil, err
	}

	return res.Children[0], res.Children[1], nil
}

func (engine *Engine) SaveAllCols(tblDbName string, table Table) error {
	for idx, col := range table.colNames {
		fqColName := tblDbName + "." + col
		colVals := table.cols[idx]
		crc := Checksum(colVals)
		metadata := ColMetdata{
			Column:         fqColName,
			Formula:        "",
			Deps:           nil,
			Checksum:       crc,
			LastModifiedMs: time.Now().UnixMilli(),
		}
		engine.metadata[fqColName] = metadata
		err := engine.loader.SaveColumn(metadata, colVals, *saveBackupsFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (engine *Engine) Evaluate(cmd string) (int, error) {
	code, _, err := engine.EvaluateCmd(cmd)
	return code, err
}

func (engine *Engine) EvaluateCmd(cmd string) (int, string, error) {
	cmd = strings.TrimSpace(cmd)
	engine.cmd = cmd
	if strings.Count(cmd, EQUAL_SEP) != 1 {
		// TODO: This doesnt handle a = CONCAT('=', b) correctly
		return ERROR_BAD_LINE, "", fmt.Errorf("0 or more than 1 = separaters found")
	}

	parts := strings.Split(cmd, EQUAL_SEP)
	_lhs, _rhs := parts[0], parts[1]

	// Load inline csv data
	matches := CSV_REGEX.FindStringSubmatch(_rhs)
	if len(matches) >= 2 {
		err := engine.LoadCsvContent(matches[1], _lhs)
		if err != nil {
			return ERROR_BAD_CSV_LINES, _lhs, err
		}
		return CODE_OK, _lhs, nil
	}

	lhsN, rhsN, err := engine.Parse(cmd)
	if err != nil {
		return GENERAL_ERROR, _lhs, err
	}
	engine.lhs = lhsN
	_, lhsTblName, lhsCol := parseTbl(_lhs)
	lhsN.TblDbName = lhsTblName

	// New dbs can be initialized with other load commands
	if lhsN.Type == TYPE_DB {
		val, err := engine.evaluateTreeInt(rhsN, 0)

		if strings.Contains(_rhs, FUNC_XRANGE) {
			if newDb, ok := val.(Table); ok {
				engine.dbs[lhsN.TblDbName] = newDb
				err := engine.SaveAllCols(_lhs, newDb)
				if err != nil {
					return ERROR_IN_PERSISTENCE, _lhs, err
				}
				return CODE_OK, _lhs, nil
			} else {
				return GENERAL_ERROR, _lhs, fmt.Errorf("bad xrange")
			}
		}

		if err != nil {
			return GENERAL_ERROR, _lhs, err
		}
		return CODE_OK, _lhs, nil
	}
	if lhsN.Type != TYPE_TBL {
		return ERROR_BAD_LINE, _lhs, fmt.Errorf("table expected in lhs: %s", lhsN.Text)
	}

	lhsN.Col = lhsCol
	if lhsN.Col != "" {
		// New column
		if table, exists := engine.tables[lhsN.TblDbName]; !exists {
			return GENERAL_ERROR, _lhs, fmt.Errorf("table doesnt exist: %s", lhsN.TblDbName)
		} else {
			if lo.Contains(table.colNames, lhsN.Col) {
				return GENERAL_ERROR, _lhs, fmt.Errorf("column already exists in table: %s", lhsN.Text)
			}

			engine.lhs = lhsN

			var colVals []ResultValueType
			metadata, cachedVals, found := engine.loader.GetColumn(_lhs)
			if found && engine.isColMetadataValid(metadata, rhsN.Text) {
				colVals = cachedVals
				engine.metadata[_lhs] = metadata
			} else {
				fmt.Println("Not found in metadata cache: ", _lhs)
				colVals = make([]ResultValueType, table.nRows)
				err := engine.evalParallel(cmd, colVals, rhsN, table)
				if err != nil {
					return GENERAL_ERROR, _lhs, err
				}

				crc := Checksum(colVals)
				terms := lo.Uniq(engine.captureTermsInsideNode(rhsN))
				deps := make([]_Dep, 0, len(terms))
				for _, term := range terms {
					if m, exists := engine.metadata[term]; !exists {
						return ERROR_IN_PERSISTENCE, _lhs, fmt.Errorf("didnt find metadata for dependency: %s", term)
					} else {
						deps = append(deps, _Dep{Term: term, Checksum: m.Checksum, LastModifiedMs: m.LastModifiedMs})
					}
				}
				metadata := ColMetdata{
					Column:         _lhs,
					Formula:        rhsN.Text,
					Deps:           deps,
					Checksum:       crc,
					LastModifiedMs: time.Now().UnixMilli(),
				}
				engine.metadata[_lhs] = metadata
				err = engine.loader.SaveColumn(metadata, colVals, *saveBackupsFlag)
				if err != nil {
					return ERROR_IN_PERSISTENCE, _lhs, err
				}
			}

			table.colNames = append(table.colNames, lhsN.Col)
			table.cols = append(table.cols, colVals)
			engine.tables[lhsN.TblDbName] = table
		}

		return CODE_OK, _lhs, nil
	} else {
		// New table declaration
		if _, exists := engine.tables[lhsN.TblDbName]; exists {
			return GENERAL_ERROR, _lhs, fmt.Errorf("table already exist: %s", lhsN.TblDbName)
		}

		if strings.Contains(_rhs, FUNC_FILTER) {
			// TODO: Fix this
			args := rhsN.Children[0].Children[1]
			noop(rhsN, args)
			_tables, err := engine.lookupTables([]string{args.Children[0].Text})
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			terms, err := engine.captureTermsInsideLookup(rhsN)
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			noop(terms)
			fcList, err := convertTerminalsToFC(terms[2:])
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			noop(fcList)

			mine, _ := engine._findMineOther(fcList[0], parseTblStruct(args.Children[0].Text))
			cols := make([][]ResultValueType, len(_tables[0].colNames))
			nRows := 0
			for i := 0; i < _tables[0].nRows; i++ {
				conds := engine._evaluateConditions(fcList, mine, i, i)
				if !lo.Contains(conds, false) {
					for idx := range _tables[0].colNames {
						cols[idx] = append(cols[idx], _tables[0].cols[idx][i])
					}
					nRows++
				} else {
					noop("didnt match")
				}
			}

			_table := Table{colNames: _tables[0].colNames, nRows: nRows, cols: cols}
			engine.tables[lhsN.TblDbName] = _table
			err = engine.SaveAllCols(_lhs, _table)
			if err != nil {
				return ERROR_IN_PERSISTENCE, _lhs, err
			}
			return CODE_OK, _lhs, nil
		}

		if strings.Contains(_rhs, FUNC_CONCAT_TABLE) {
			terms, err := engine.captureTermsInsideLookup(rhsN)
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			if len(terms) == 0 || terms[0] != FUNC_CONCAT_TABLE {
				return GENERAL_ERROR, _lhs, fmt.Errorf("wrong syntax for CONCAT_TABLE")
			}
			_table, err := engine.concatTables(terms[1:])
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			engine.tables[lhsN.TblDbName] = _table
			err = engine.SaveAllCols(_lhs, _table)
			if err != nil {
				return ERROR_IN_PERSISTENCE, _lhs, err
			}
			return CODE_OK, _lhs, nil
		}

		if strings.Contains(_rhs, FUNC_RETAIN_COLS) {
			terms, err := engine.captureTermsInsideLookup(rhsN)
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			if len(terms) != 3 || terms[0] != FUNC_RETAIN_COLS {
				return GENERAL_ERROR, _lhs, fmt.Errorf("wrong syntax for RETAIN_COLS")
			}
			tables, err := engine.lookupTables(terms[1:2])
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			colNames := strings.Split(terms[2], ",")
			cols := make([][]ResultValueType, len(colNames))
			for jdx, col := range colNames {
				idx := lo.IndexOf(tables[0].colNames, col)
				if idx < 0 {
					return GENERAL_ERROR, _lhs, fmt.Errorf("col not found: %s %s", col, terms[1])
				}
				cols[jdx] = tables[0].cols[idx]
			}
			table := Table{colNames: colNames, nRows: tables[0].nRows, cols: cols}
			engine.tables[lhsN.TblDbName] = table
			err = engine.SaveAllCols(_lhs, table)
			if err != nil {
				return ERROR_IN_PERSISTENCE, _lhs, err
			}
			return CODE_OK, _lhs, nil
		}

		if strings.Contains(_rhs, FUNC_CARTESIAN) {
			// Cartesian
			terms, err := engine.captureTermsInsideLookup(rhsN)
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}
			terms = lo.Filter(terms, func(item string, index int) bool { return item != FUNC_CARTESIAN })
			colNames, cols, nRows, err := engine.constructCartesianTable(rhsN.Text, terms)
			if err != nil {
				return GENERAL_ERROR, _lhs, err
			}

			_table := Table{colNames: colNames, cols: cols, nRows: nRows}
			engine.tables[lhsN.TblDbName] = _table
			err = engine.SaveAllCols(_lhs, _table)
			if err != nil {
				return ERROR_IN_PERSISTENCE, _lhs, err
			}
			return CODE_OK, _lhs, nil
		}

		// tbl.X = db.Y declaration
		_type, _name, _col := parseTbl(_rhs)
		if _type != TYPE_DB || _col != "" {
			return GENERAL_ERROR, _lhs, fmt.Errorf("table can only be initialized from existing dbs: %s", cmd)
		}
		_table := engine.dbs[_name]
		engine.tables[lhsN.TblDbName] = _table
		err = engine.SaveAllCols(_lhs, _table)
		if err != nil {
			return ERROR_IN_PERSISTENCE, _lhs, err
		}
		return CODE_OK, _lhs, nil
	}
}

func findFilterNode(node *Node) *Node {
	if node.Ctx == CTX_FunctionContext && len(node.Children) == 2 && node.Children[0].Type == TYPE_CONTEXT_FN && node.Children[0].Terminal == FUNC_LOOKUP {
		return node
	}
	for _, child := range node.Children {
		f := findFilterNode(child)
		if f != nil {
			return f
		}
	}
	return nil
}

// NOTE: Evaluation can be parallelized here, because its guaranteed that each row of this column can be evaluated independently
func (engine *Engine) evalParallel(cmd string, colVals []ResultValueType, rhsN *Node, table Table) error {
	nRowsPerBatch := engine.nRowsPerLookupBatch
	nBatches := 1
	if strings.Contains(rhsN.Text, FUNC_LOOKUP) && table.nRows > nRowsPerBatch {
		nBatches = int(math.Ceil(float64(table.nRows) / float64(nRowsPerBatch)))
	}

	// Preprocess
	fNode := findFilterNode(rhsN)
	if fNode == nil && strings.Contains(rhsN.Text, FUNC_LOOKUP) {
		panic("Filter node expected")
	}
	if fNode != nil {
		fCtx, err := engine.evalPlanForLookupFn(fNode)
		if err != nil {
			return err
		}
		engine.fCtx = fCtx
	}

	if nBatches == 1 {
		for i := 0; i < table.nRows; i++ {
			engine.lhsIndex = i
			val, err := engine.evaluateTreeInt(rhsN, 0)
			if err != nil {
				return fmt.Errorf("error executing cmd at row: %d %s %s", i, cmd, err)
			}
			colVals[i] = val
		}
		return nil
	}

	// Create index on column if it doesnt exist
	otherCol := engine.fCtx.otherCol.fqName
	if _, _found := engine.indexes[otherCol]; !_found {
		_index2, _num := engine.createColIndex(otherCol, 0)
		_index := ColIndex{index: _index2, tbl: otherCol, numRowsSynced: _num}
		engine.indexes[otherCol] = _index
	}

	errors := make([]error, nBatches)
	var wg sync.WaitGroup
	wg.Add(nBatches)
	for j := 0; j < nBatches; j++ {
		go func(engine *Engine, idx int) {
			eng2 := &Engine{
				app:      engine.app,
				dbs:      engine.dbs,
				tables:   engine.tables,
				cmdIdx:   engine.cmdIdx,
				cmd:      engine.cmd,
				indexes:  engine.indexes,
				lhs:      engine.lhs,
				lhsIndex: -1,
				fCtx:     engine.fCtx,
			}
			startIdx := nRowsPerBatch * idx
			endIdx := MinV(nRowsPerBatch*(idx+1), table.nRows)
			for i := startIdx; i < endIdx; i++ {
				eng2.lhsIndex = i
				val, err := eng2.evaluateTreeInt(rhsN, 0)
				if err != nil {
					errStr := fmt.Errorf("error executing cmd at row: %d %s %s", i, cmd, err)
					fmt.Println(errStr)
					errors[idx] = errStr
					return
				}
				colVals[i] = val
			}

			wg.Done()
		}(engine, j)
	}
	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) evaluateTreeInt(node *Node, level int) (ResultValueType, error) {
	if node.Ctx == CTX_ExpressionContext {
		if len(node.Children) == 1 && (node.Children[0].Ctx == CTX_FunctionContext || node.Children[0].Ctx == CTX_FunctionNoArgsContext || node.Children[0].Ctx == CTX_TerminalNode || node.Children[0].Ctx == CTX_ExpressionContext) {
			val, err := engine.evaluateTreeInt(node.Children[0], level+1)
			return val, err
		}
		if len(node.Children) == 1 && node.Children[0].Ctx == CTX_BoolContext {
			val, err := engine.evaluateTreeInt(node.Children[0].Children[0], level+1)
			return val, err
		}
	}

	// a operator b
	if (node.Ctx == CTX_ExpressionContext || node.Ctx == CTX_FunctionContext) && len(node.Children) == 3 && (node.Children[1].Ctx == CTX_MathOpContext || node.Children[1].Ctx == CTX_SetOpContext || node.Children[1].Ctx == CTX_LogicalFnContext) {
		aVal, err := engine.evaluateTreeInt(node.Children[0], level+1)
		if err != nil {
			return nil, err
		}
		bVal, err := engine.evaluateTreeInt(node.Children[2], level+1)
		if err != nil {
			return nil, err
		}

		fn := node.Children[1].Text
		if _, ok := ALL_FUNCTIONS[fn]; !ok {
			return nil, fmt.Errorf("bad function: %s", fn)
		}
		_func := ALL_FUNCTIONS[fn]

		val, err := _func(engine.evaluatorCtx, aVal, bVal)
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	// function
	cond1 := node.Ctx == CTX_ExpressionContext && len(node.Children) == 1 && node.Children[0].Ctx == CTX_FunctionNoArgsContext
	cond2 := node.Ctx == CTX_ExpressionContext && len(node.Children) == 2 && node.Children[0].Ctx == CTX_FunctionContext && node.Children[1].Ctx == CTX_ArgumentsContext
	cond3 := node.Ctx == CTX_ExpressionContext && len(node.Children) > 0 && (node.Children[0].Ctx == CTX_FunctionNoArgsContext || node.Children[0].Ctx == CTX_FunctionContext)
	cond4 := node.Ctx == CTX_FunctionContext && len(node.Children) == 2 && (node.Children[0].Ctx == CTX_TerminalNode && node.Children[1].Ctx == CTX_ArgumentsContext)
	cond5 := node.Ctx == CTX_FunctionNoArgsContext && len(node.Children) == 1 && node.Children[0].Ctx == CTX_TerminalNode
	if cond1 || cond2 || cond3 || cond4 || cond5 {
		fn := node.Children[0].Text
		if node.Children[0].Ctx == CTX_FunctionNoArgsContext {
			// Trim brackets from fn
			fn = strings.ReplaceAll(strings.ReplaceAll(fn, "(", ""), ")", "")
		}
		if _, ok := ALL_FUNCTIONS[fn]; !ok {
			return nil, fmt.Errorf("bad function: %s", fn)
		}
		if fn == FUNC_LOOKUP {
			res, err := engine.processLookupFn(node, engine.fCtx, level)
			return res, err
		}

		argVals := make([]ResultValueType, 0)
		if len(node.Children) > 1 {
			for _, child := range node.Children[1].Children {
				val, err := engine.evaluateTreeInt(child, level+1)
				if err != nil {
					return nil, err
				}
				argVals = append(argVals, val)
			}
		}

		if fn == FUNC_LOAD_CSV { // TODO: Move to functions.go
			fName := interfaceToString(argVals[0])
			err := engine.LoadCsvFile(fName, engine.lhs.Text, -1)
			if err != nil {
				return GENERAL_ERROR, err
			}
			return CODE_OK, nil
		}

		_func := ALL_FUNCTIONS[fn]
		val, err := _func(engine.evaluatorCtx, argVals)
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	// terminal node
	if node.Ctx == CTX_TerminalNode {
		switch node.Type {
		case TYPE_TBL:
			val, err := engine.resolveTblDb(node)
			return val, err
		case TYPE_DB:
			val, err := engine.resolveTblDb(node)
			return val, err
		case TYPE_STRING:
			return node.Terminal, nil
		case TYPE_BOOLEAN:
			return node.Terminal, nil
		case TYPE_NUMBER:
			return node.Terminal, nil
		default:
			return nil, fmt.Errorf("unexpected terminal node: %s", node.Text)
		}
	}

	return nil, fmt.Errorf("dont know how to evaluate node: %s %s", node.Ctx, node.Text)
}

func (engine *Engine) resolveTblDb(node *Node) (ResultValueType, error) {
	var _map map[string]Table
	if node.Type == TYPE_DB {
		_map = engine.dbs
	} else if node.Type == TYPE_TBL {
		_map = engine.tables
	} else {
		return nil, fmt.Errorf("cant lookup table or db: %v", node)
	}

	if table, exists := _map[node.TblDbName]; !exists {
		return nil, fmt.Errorf("table or db doesnt exist: %s %s", node.Type, node.TblDbName)
	} else {
		idx := lo.IndexOf(table.colNames, node.Col)
		if idx < 0 {
			return nil, fmt.Errorf("col doesnt exist in table or db: %s %s %s", node.Type, node.TblDbName, node.Col)
		}
		return table.cols[idx][engine.lhsIndex], nil
	}
}

func (engine *Engine) dfsParse(t antlr.Tree) (Node, error) {
	numChildren := t.GetChildCount()
	var ctx string = ""
	var text string = ""

	switch tt := t.(type) {
	case *parser.ProgContext:
		ctx = CTX_ProgContext
		text = tt.GetText()
	case *parser.ExpressionContext:
		ctx = CTX_ExpressionContext
		text = tt.GetText()
	case *parser.FunctionContext:
		ctx = CTX_FunctionContext
		text = tt.GetText()
	case *parser.FunctionNoArgsContext:
		ctx = CTX_FunctionNoArgsContext
		text = tt.GetText()
	case *parser.MathFnContext:
		ctx = CTX_MathOpContext
		text = tt.GetText()
	case *parser.SetFnContext:
		ctx = CTX_SetOpContext
		text = tt.GetText()
	case *parser.LogicalFnContext:
		ctx = CTX_LogicalFnContext
		text = tt.GetText()

	case *parser.ArgumentsContext:
		ctx = CTX_ArgumentsContext
		text = tt.GetText()
	case *parser.BoolContext:
		ctx = CTX_BoolContext
		text = tt.GetText()

	// Lhs
	case *parser.LhsContext:
		ctx = CTX_LhsContext
		text = tt.GetText()
		_type, _, err := typeOfTerminalNode(ctx, text)
		if err != nil {
			return Node{}, err
		}
		return Node{Ctx: ctx, Type: _type, Text: text}, nil

	// Error
	case antlr.ErrorNode:
		ctx = CTX_Error
		return Node{Ctx: ctx, Type: "", Text: ""}, nil

	// Terminal node
	case antlr.TerminalNode:
		ctx = CTX_TerminalNode
		text = tt.GetText()
		_type, _res, err := typeOfTerminalNode(ctx, text)
		if err != nil {
			return Node{}, err
		}
		node := Node{Ctx: ctx, Type: _type, Text: text, Terminal: _res}
		if _type == TYPE_TBL || _type == TYPE_DB {
			_, _s1, _s2 := parseTbl(text)
			node.TblDbName = _s1
			node.Col = _s2
		}
		return node, nil

	default:
		return Node{}, fmt.Errorf("unexpected tree node: %s %s", ctx, text)
	}
	noop(ctx, text)

	childs := make([]*Node, 0)
	for i := 0; i < numChildren; i++ {
		child := t.GetChild(i)
		res, err := engine.dfsParse(child)
		if err != nil {
			return Node{}, err
		}
		if res.Type != TYPE_IGNORE {
			childs = append(childs, &res)
		}
	}
	return Node{Ctx: ctx, Type: "", Text: text, Children: childs}, nil
}

func typeOfTerminalNode(ctx, text string) (string, ResultValueType, error) {
	if strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'") {
		// string
		return TYPE_STRING, text[1 : len(text)-1], nil
	} else if text == "OR" {
		return TYPE_MATH_FN, text, nil
	} else if text == "AND" {
		return TYPE_MATH_FN, text, nil
	} else if text == "TRUE" {
		// boolean
		return TYPE_BOOLEAN, true, nil
	} else if text == "FALSE" {
		// boolean
		return TYPE_BOOLEAN, false, nil
	} else if text == "(" || text == ")" || text == "=" || text == "," || text == "<EOF>" {
		// ignore
		return TYPE_IGNORE, text, nil
	} else if len(text) > 0 && (strings.Contains(DIGITS, string(text[0])) || (len(text) >= 2 && text[0] == '-' && strings.Contains(DIGITS, string(text[1])))) {
		// number
		v, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse float: %s %s", text, err)
		}
		return TYPE_NUMBER, v, nil
	} else if _, ok := MATH_OPERATOR_FUNCTIONS[text]; ok {
		// math function
		return TYPE_MATH_OPERATOR_FN, text, nil
	} else if _, ok := MATH_FUNCTIONS[text]; ok {
		return TYPE_MATH_FN, text, nil
	} else if _, ok := SET_FUNCTIONS[text]; ok {
		// set function
		return TYPE_SET_FN, text, nil
	} else if _, ok := STRING_FUNCTIONS[text]; ok {
		// string function
		return TYPE_STRING_FN, text, nil
	} else if _, ok := CONTEXT_FUNCTIONS[text]; ok {
		// context function
		return TYPE_CONTEXT_FN, text, nil
	} else if strings.HasPrefix(text, TYPE_TBL) {
		// table term
		return TYPE_TBL, text, nil
	} else if strings.HasPrefix(text, TYPE_CONST) {
		// constant
		return TYPE_CONST, text, nil
	} else if strings.HasPrefix(text, TYPE_STAT) {
		// stat
		return TYPE_STAT, text, nil
	} else if strings.HasPrefix(text, TYPE_DB) {
		// db term
		return TYPE_DB, text, nil
	} else {
		return "", nil, fmt.Errorf("unexpected terminal node: %s %s", ctx, text)
	}
}

func (engine *Engine) getTableOrDb(lhs string, numRows int) (string, error) {
	_type, _name, _ := parseTbl(lhs)

	if _type == TYPE_TBL {
		if table, found := engine.tables[_name]; !found {
			return "", fmt.Errorf("table doesnt exist: %s", _name)
		} else {
			return engine._getTableOrDb(table, numRows)
		}
	} else if _type == TYPE_DB {
		if table, found := engine.dbs[_name]; !found {
			return "", fmt.Errorf("db doesnt exist: %s", _name)
		} else {
			return engine._getTableOrDb(table, numRows)
		}
	} else {
		return "", fmt.Errorf("bad type: %s", _type)
	}
}

func (engine *Engine) _getTableOrDb(table Table, numRows int) (string, error) {
	nr := lo.Min([]int{numRows, table.nRows})
	cols := make([][]string, len(table.colNames))
	for i := 0; i < len(cols); i++ {
		cols[i] = lo.Map(table.cols[i][0:nr], func(item ResultValueType, index int) string { return interfaceToString(item) })
	}

	out := ReturnTable{ColNames: table.colNames, Cols: cols}
	_bytes, err := json.Marshal(out)
	return string(_bytes), err
}

func (engine *Engine) LoadCsvFile(file, lhs string, numLines int) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return engine.LoadCsvContent(string(data), lhs)
}

func (engine *Engine) LoadCsvContent(content, lhs string) error {
	csvReader := csv.NewReader(strings.NewReader(content))
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("unable to parse csv: %s", err)
	}

	// Table / db shouldnt already exist
	_type, _name, _ := parseTbl(lhs)
	if _, found := engine.tables[_name]; _type == TYPE_TBL && found {
		return fmt.Errorf("table already exists: %s", _name)
	} else if _, found := engine.dbs[_name]; _type == TYPE_DB && found {
		return fmt.Errorf("db already exists: %s", _name)
	}

	// Create table in columnar format
	table := Table{colNames: records[0], nRows: 0, cols: make([][]ResultValueType, 0, len(records[0]))}
	for i := range table.colNames {
		noop(i)
		table.cols = append(table.cols, make([]ResultValueType, 0, len(records)-1))
	}
	for _, line := range records[1:] {
		table.nRows++
		for idx, v := range line {
			table.cols[idx] = append(table.cols[idx], v)
		}
	}

	_type, tableName, _ := parseTbl(lhs)
	noop(_type, tableName)
	if strings.HasPrefix(lhs, TYPE_TBL+".") {
		engine.tables[strings.Split(lhs, ".")[1]] = table
	} else if strings.HasPrefix(lhs, TYPE_DB+".") {
		engine.dbs[strings.Split(lhs, ".")[1]] = table
	} else {
		return fmt.Errorf("bad lhs: %s", lhs)
	}
	engine.SaveAllCols(lhs, table)

	return nil
}

func (engine *Engine) writeTableToFile(tblName, filename string) error {
	table, _found := engine.tables[tblName]
	if !_found {
		return fmt.Errorf("table not found: %s", tblName)
	}

	f2, err := os.OpenFile(filename, FILE_FLAGS, 0644)
	if err != nil {
		return err
	}
	w2 := csv.NewWriter(f2)
	w2.Write(table.colNames)
	for i := 0; i < table.nRows; i++ {
		row := make([]string, 0)
		for _, col := range table.cols {
			row = append(row, interfaceToString(col[i]))
		}
		w2.Write(row)
	}

	w2.Flush()
	defer f2.Close()
	return nil
}
