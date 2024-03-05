package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/btree"
)

func insertValIntoIndex(_index *btree.BTreeG[ColIdxItem], val string, idx int) bool {
	if sItem, _found := _index.Get(ColIdxItem{val, nil}); _found {
		sItem.idxList = append(sItem.idxList, idx)
		_index.Set(sItem)
		return true
	} else {
		_index.Set(ColIdxItem{val, []int{idx}})
		return false
	}
}

func verifyFilterTerms(cols []string, lhs string) bool {
	mm := make(map[string]bool)
	for _, c := range cols {
		_type, _name, _ := parseTbl(c)
		if _type != TYPE_TBL && _type != TYPE_DB {
			continue
		}
		key := _type + "." + _name
		mm[key] = true
	}
	if len(mm) > 2 {
		return false
	}
	if len(mm) == 1 {
		return false
	}

	_t, _n, _ := parseTbl(lhs)
	if _, found := mm[_t+"."+_n]; !found {
		return false
	}

	return true
}

func convertTerminalsToFC(terminals []string) ([]_FilterCondition, error) {
	if len(terminals)%3 != 0 {
		return nil, fmt.Errorf("cant parse terminals to filter conditions: bad len %d", len(terminals))
	}

	list := make([]_FilterCondition, 0)
	for i := 0; i < len(terminals); i += 3 {
		list = append(list, _FilterCondition{col1: parseTblStruct(terminals[i]), op: terminals[i+1], col2: parseTblStruct(terminals[i+2])})
	}
	return list, nil
}

func parseTbl(tbl string) (string, string, string) {
	splits := strings.Split(tbl, ".")
	if len(splits) < 2 || len(splits) > 3 {
		return "", "", ""
	}
	t := TYPE_TBL
	if strings.HasPrefix(tbl, "db") {
		t = TYPE_DB
	}
	if len(splits) == 2 {
		return t, splits[1], ""
	}
	return t, splits[1], splits[2]
}

func parseTblStruct(tbl string) _Tbl {
	_type, _name, _col := parseTbl(tbl)
	return _Tbl{Type: _type, Name: _name, Col: _col, fqName: tbl}
}

func itemLessFn(a, b ColIdxItem) bool {
	ret := a.val < b.val
	return ret
}

const (
	CTX_ProgContext       = "ProgContext"
	CTX_LhsContext        = "LhsContext"
	CTX_ExpressionContext = "ExpressionContext"
	CTX_FunctionContext   = "FunctionContext"
	CTX_MathOpContext     = "MathFnContext"
	CTX_SetOpContext      = "SetFnContext"
	CTX_LogicalFnContext  = "LogicalFnContext"
	CTX_ArgumentsContext  = "ArgumentsContext"
	CTX_BoolContext       = "BoolContext"
	CTX_TerminalNode      = "TerminalNode"
	CTX_Error             = "ErrorContext"
)

const (
	TYPE_IGNORE           = "ignore"
	TYPE_STRING           = "string"
	TYPE_BOOLEAN          = "boolean"
	TYPE_NUMBER           = "number"
	TYPE_MATH_OPERATOR_FN = "mathOperator"
	TYPE_MATH_FN          = "mathFn"
	TYPE_SET_FN           = "setOperator"
	TYPE_STRING_FN        = "stringFn"
	TYPE_CONTEXT_FN       = "contextFn"
	TYPE_TBL              = "tbl"
	TYPE_DB               = "db"
	TYPE_STAT             = "stat"
	TYPE_CONST            = "const"
	TYPE_ARGS             = "args"
)
const DIGITS = "0123456789"

var CSV_REGEX = regexp.MustCompile(`(?sm)^CSV ?\( ?\x60(.*)\x60 ?\)$`)
var DB_REGEX = regexp.MustCompile(`(?sm)^db.([_a-zA-Z][_a-zA-Z0-9]*)$`)

const EQUAL_SEP = " = "
const (
	CODE_OK              = 0
	GENERAL_ERROR        = 1
	ERROR_BAD_LINE       = 2
	ERROR_BAD_CSV_LINES  = 3
	ERROR_IN_PERSISTENCE = 4
)

type Table struct {
	colNames []string
	nRows    int
	cols     [][]ResultValueType
}

type ColIndex struct {
	index         *btree.BTreeG[ColIdxItem]
	tbl           string
	numRowsSynced int
}

type ColIdxItem struct {
	val     string
	idxList []int
}

type ResultValueType any

type Result struct {
	Ctx  string
	Type string
	Val  ResultValueType
	Ref  string
}

type _FilterCondition struct {
	col1 _Tbl
	op   string
	col2 _Tbl
}

type _Tbl struct {
	Type   string
	Name   string
	Col    string
	fqName string
}

type ReturnTable struct {
	ColNames []string   `json:"colNames"`
	Cols     [][]string `json:"cols"`
}
