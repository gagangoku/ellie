package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samber/lo"
	"github.com/tidwall/btree"
)

type FilterContext struct {
	terms         []string
	fcList        []_FilterCondition
	otherColToRet _Tbl
	myCol         _Tbl
	otherCol      _Tbl
	lhsTbl        _Tbl
}

func (engine *Engine) evalPlanForLookupFn(node *Node) (*FilterContext, error) {
	child := node.Children[0]
	if child.Type != TYPE_CONTEXT_FN || child.Terminal != FUNC_LOOKUP {
		return nil, fmt.Errorf("unexpected 2")
	}

	terms, err := engine.captureTermsInsideLookup(node)
	if err != nil {
		return nil, err
	}
	if len(terms) < 5 || terms[0] != FUNC_LOOKUP {
		return nil, fmt.Errorf("LOOKUP syntax: LOOKUP(fully_qualified_col, <eq condition1>, <condition>*)")
	}
	terms = terms[1:]
	isValid := verifyFilterTerms(terms, engine.lhs.Text)
	if !isValid {
		return nil, fmt.Errorf("LOOKUP function only supports 1 other table")
	}

	otherColToRet := terms[0]
	fcList, err := convertTerminalsToFC(terms[1:])
	if err != nil {
		return nil, err
	}
	noop(otherColToRet, fcList)

	if fcList[0].op != FUNC_EQ {
		return nil, fmt.Errorf("first condition of LOOKUP should be equality check")
	}
	myCol, otherCol := engine._findMineOther(fcList[0], parseTblStruct(engine.lhs.Text))
	if otherCol.fqName == "" {
		return nil, fmt.Errorf("otherCol not found: %s", node.Text)
	}

	return &FilterContext{
		terms:         terms,
		fcList:        fcList,
		otherColToRet: parseTblStruct(otherColToRet),
		otherCol:      otherCol,
		myCol:         myCol,
		lhsTbl:        parseTblStruct(engine.lhs.Text),
	}, nil
}

func (engine *Engine) processLookupFn(node *Node, fCtx *FilterContext, level int) (ResultValueType, error) {
	if fCtx == nil {
		panic("unexpected")
	}
	res, err := engine.resolveTblDb(&Node{Type: fCtx.myCol.Type, TblDbName: fCtx.myCol.Name, Col: fCtx.myCol.Col})
	if err != nil {
		return nil, err
	}

	otherCol := fCtx.otherCol.fqName
	_index, _found := engine.indexes[otherCol]
	if !_found {
		// create index on column
		_index2, _num := engine.createColIndex(otherCol, 0)
		_index = ColIndex{index: _index2, tbl: otherCol, numRowsSynced: _num}
		engine.indexes[otherCol] = _index
	}
	lookup, _found := _index.index.Get(ColIdxItem{interfaceToString(res), nil})
	if !_found {
		// Didnt match anything in the other table
		return []ResultValueType{}, nil
	}

	// Found corresponding entries in other table
	vals := make([]ResultValueType, 0, len(lookup.idxList))
	for _, otherIdx := range lookup.idxList {
		conds := engine._evaluateConditions(fCtx.fcList[1:], fCtx.lhsTbl, engine.lhsIndex, otherIdx)
		if lo.Contains(conds, false) {
			continue
		}

		v := engine._getValue(fCtx.otherColToRet, otherIdx)
		vals = append(vals, v)
	}

	return vals, nil
}

func (s *Engine) createColIndex(tbl string, prevIdx int) (*btree.BTreeG[ColIdxItem], int) {
	_type, tableName, colName := parseTbl(tbl)
	noop(_type, tableName, colName)
	var table Table
	if _type == TYPE_DB {
		table = s.dbs[tableName]
	} else {
		table = s.tables[tableName]
	}

	if prevIdx < 0 {
		prevIdx = 0
	}

	tree := btree.NewBTreeG[ColIdxItem](itemLessFn)
	idx := lo.IndexOf(table.colNames, colName)
	for idx, v := range table.cols[idx][prevIdx:] {
		insertValIntoIndex(tree, interfaceToString(v), idx)
	}
	return tree, len(table.cols[idx])
}

func (e *Engine) _evaluateConditions(fcList []_FilterCondition, lhs _Tbl, lhsIdx, rhsIdx int) []bool {
	conds := make([]bool, 0)
	for _, fc := range fcList {
		mine, other := e._findMineOther(fc, lhs)
		v1 := e._getValue(mine, lhsIdx)
		v2 := e._getValue(other, rhsIdx)

		_func := ALL_FUNCTIONS[fc.op]

		var res ResultValueType
		var err error
		if fc.col1 == mine {
			res, err = _func(nil, v1, v2)
		} else {
			res, err = _func(nil, v2, v1)
		}
		if err != nil {
			noop(err)
		}
		c := interfaceToBool(res)
		conds = append(conds, c)
	}
	return conds
}

func (e *Engine) _getValue(col _Tbl, idx int) ResultValueType {
	_t, _n, _c := col.Type, col.Name, col.Col
	if _t == TYPE_DB {
		t := e.dbs[_n]
		_idx := lo.IndexOf(t.colNames, _c)
		return t.cols[_idx][idx]
	} else if _t == TYPE_TBL {
		t := e.tables[_n]
		_idx := lo.IndexOf(t.colNames, _c)
		return t.cols[_idx][idx]
	} else if strings.HasPrefix(col.fqName, "'") && strings.HasSuffix(col.fqName, "'") {
		// A string
		ret := col.fqName[1 : len(col.fqName)-1]
		return ret
	} else {
		// Return as is
		return col.fqName
	}
}

func (e *Engine) _findMineOther(fc _FilterCondition, lhs _Tbl) (_Tbl, _Tbl) {
	_t, _n := lhs.Type, lhs.Name
	_t2, _n2 := fc.col1.Type, fc.col1.Name
	if _t == _t2 && _n == _n2 {
		return fc.col1, fc.col2
	}
	return fc.col2, fc.col1
}

func (engine *Engine) captureTermsInsideLookup(t *Node) ([]string, error) {
	if t.Ctx == CTX_TerminalNode && t.Type != TYPE_IGNORE {
		return []string{interfaceToString(t.Terminal)}, nil
	}

	list := make([]string, 0)
	for _, child := range t.Children {
		res, err := engine.captureTermsInsideLookup(child)
		if err != nil {
			return nil, err
		}
		list = append(list, res...)
	}
	return list, nil
}

func (engine *Engine) captureTermsInsideNode(t *Node) []string {
	if t.Ctx == CTX_TerminalNode && (t.Type == TYPE_DB || t.Type == TYPE_TBL) {
		return []string{t.Text}
	}

	list := make([]string, 0)
	for _, child := range t.Children {
		res := engine.captureTermsInsideNode(child)
		list = append(list, res...)
	}
	return list
}

func (engine *Engine) constructCartesianTable(cmd string, terms []string) ([]string, [][]ResultValueType, int, error) {
	type _Parsed struct {
		tbl   _Tbl
		table Table
	}
	parsedTerms := make([]_Parsed, 0, len(terms))
	for _, term := range terms {
		tblTerm := parseTblStruct(term)
		_ns, _name, _col := tblTerm.Type, tblTerm.Name, tblTerm.Col
		if _ns == "" || _name == "" || _col == "" {
			return nil, nil, 0, fmt.Errorf("bad term in cartesian: %s %s", term, cmd)
		}

		var table Table
		if _ns == TYPE_TBL {
			t, ok := engine.tables[_name]
			if !ok {
				return nil, nil, 0, fmt.Errorf("table not found: %s %s", term, cmd)
			}
			table = t
		} else if _ns == TYPE_DB {
			t, ok := engine.dbs[_name]
			if !ok {
				return nil, nil, 0, fmt.Errorf("db not found: %s %s", term, cmd)
			}
			table = t
		} else {
			return nil, nil, 0, fmt.Errorf("bad term in cartesian: %s %s", term, cmd)
		}

		idx := lo.IndexOf(table.colNames, _col)
		if idx < 0 {
			return nil, nil, 0, fmt.Errorf("col not found: %s %s", term, cmd)
		}
		parsedTerms = append(parsedTerms, _Parsed{tbl: tblTerm, table: table})
	}

	_colNames := lo.Map(parsedTerms, func(item _Parsed, index int) string { return item.tbl.Col })
	if len(lo.Uniq(_colNames)) != len(_colNames) {
		return nil, nil, 0, fmt.Errorf("duplicate column in CARTESIAN not supported: %s", cmd)
	}

	// Group by table name
	grouped := lo.GroupBy(parsedTerms, func(item _Parsed) string { return item.tbl.Type + "." + item.tbl.Name })
	groupNames := lo.Keys(grouped)
	sort.Strings(groupNames)
	combLens := make([]int, 0, len(groupNames))
	colValues := make([][]ResultValueType, 0, len(groupNames))

	colNames := make([]string, 0, len(terms))
	for _, _g := range groupNames {
		group := grouped[_g]
		for _, parsed := range group {
			colNames = append(colNames, parsed.tbl.Col)
		}
	}

	groupToUniqVals := make([][][]ResultValueType, 0, len(combLens))
	for _, _g := range groupNames {
		group := grouped[_g]
		mmap := make(map[string]bool)
		vList := make([][]ResultValueType, 0)

		nRows := group[0].table.nRows
		for rowIdx := 0; rowIdx < nRows; rowIdx++ {
			row := make([]ResultValueType, 0, len(group))
			var sb strings.Builder
			for _, g := range group {
				colIdx := lo.IndexOf(g.table.colNames, g.tbl.Col)
				if colIdx < 0 {
					panic("bad condition")
				}
				val := g.table.cols[colIdx][rowIdx]
				row = append(row, val)
				sb.Write([]byte(interfaceToString(val)))
				sb.Write([]byte(","))
			}
			if _, exists := mmap[sb.String()]; !exists {
				mmap[sb.String()] = true
				vList = append(vList, row)
			}
		}
		if len(mmap) == 0 {
			return colNames, lo.Map(colNames, func(item string, index int) []ResultValueType { return nil }), 0, nil
		}
		groupToUniqVals = append(groupToUniqVals, vList)
		combLens = append(combLens, len(mmap))

		for j := range group {
			noop(j)
			colValues = append(colValues, make([]ResultValueType, 0, len(mmap)))
		}
	}

	cgen, comb := NewCombinator(combLens, true)
	combIdx := 0
	for cgen.Next() {
		cgen.Product(comb)
		combIdx++

		colIdx := 0
		for i, c := range comb {
			_tmp := groupToUniqVals[i][c]
			for _, val := range _tmp {
				colValues[colIdx] = append(colValues[colIdx], val)
				colIdx++
			}
		}
	}
	return colNames, colValues, combIdx, nil
}

func (engine *Engine) lookupTables(terms []string) ([]Table, error) {
	tables := make([]Table, 0, len(terms))
	for _, term := range terms {
		_type, _name, _col := parseTbl(term)
		if _col != "" {
			return nil, fmt.Errorf("term should only have table or db: %s", term)
		}

		var tb Table
		if _type == TYPE_DB {
			if _val, exists := engine.dbs[_name]; !exists {
				return nil, fmt.Errorf("db doesnt exist: %s", _name)
			} else {
				tb = _val
			}
		} else if _type == TYPE_TBL {
			if _val, exists := engine.tables[_name]; !exists {
				return nil, fmt.Errorf("table doesnt exist: %s", _name)
			} else {
				tb = _val
			}
		} else {
			return nil, fmt.Errorf("bad term: %s", term)
		}
		tables = append(tables, tb)
	}
	return tables, nil
}

func (engine *Engine) concatTables(terms []string) (Table, error) {
	if len(terms) < 2 {
		return Table{}, fmt.Errorf("need atleast 2 tables to concat")
	}
	tables, err := engine.lookupTables(terms)
	if err != nil {
		return Table{}, err
	}
	colNames := tables[0].colNames
	if len(colNames) == 0 {
		return Table{}, fmt.Errorf("need atleast 1 column in first table to concat")
	}
	cols := make([][]ResultValueType, len(colNames))
	for jdx, col := range colNames {
		for i, table := range tables {
			idx := lo.IndexOf(table.colNames, col)
			if idx < 0 {
				return Table{}, fmt.Errorf("col not found: %s %s", col, terms[i])
			}
			cols[jdx] = append(cols[jdx], table.cols[idx]...)
		}
	}
	return Table{colNames: colNames, nRows: len(cols[0]), cols: cols}, nil
}
