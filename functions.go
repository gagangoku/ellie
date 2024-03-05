package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/samber/lo"
)

const (
	FUNC_PLUS            = "+"
	FUNC_MINUS           = "-"
	FUNC_MULT            = "*"
	FUNC_DIVIDE          = "/"
	FUNC_POW             = "^"
	FUNC_GT              = ">"
	FUNC_GE              = ">="
	FUNC_LT              = "<"
	FUNC_LE              = "<="
	FUNC_EQ              = "=="
	FUNC_NE              = "!="
	FUNC_OR              = "OR"
	FUNC_AND             = "AND"
	FUNC_SET_IN          = "IN"
	FUNC_SET_CONTAINS    = "CONTAINS"
	FUNC_CONCAT          = "CONCAT"
	FUNC_JOIN            = "JOIN"
	FUNC_STR_LEN         = "STR_LEN"
	FUNC_STR_CONTAINS    = "STR_CONTAINS"
	FUNC_STR_STARTS_WITH = "STARTS_WITH"
	FUNC_STR_ENDS_WITH   = "ENDS_WITH"
	FUNC_IF              = "IF"
	FUNC_ROW             = "ROW"
	FUNC_LOOKUP          = "LOOKUP"
	FUNC_FILTER          = "FILTER"
	FUNC_SUM             = "SUM"
	FUNC_AVERAGE         = "AVERAGE"
	FUNC_FIRST           = "FIRST"
	FUNC_LAST            = "LAST"
	FUNC_CUMULATIVE_SUM  = "CUMULATIVE_SUM"
	FUNC_MIN             = "MIN"
	FUNC_MAX             = "MAX"
	FUNC_MAXV            = "MAXV"
	FUNC_MINV            = "MINV"
	FUNC_COUNT           = "COUNT"
	FUNC_COUNT_UNIQUE    = "COUNT_UNIQUE"
	FUNC_EPOCH_MS        = "EPOCH_MS"
	FUNC_NOW_MS          = "NOW_MS"
	FUNC_CEIL            = "CEIL"
	FUNC_FLOOR           = "FLOOR"
	FUNC_ROUND           = "ROUND"
	FUNC_SYNC_DB_INC     = "SYNC_DB_INC"
	FUNC_SYNC_DB_FULL    = "SYNC_DB_FULL"
	FUNC_LOAD_CSV        = "LOAD_CSV"
	FUNC_SPLIT           = "SPLIT"
	FUNC_CARTESIAN       = "CARTESIAN"
	FUNC_XRANGE          = "XRANGE"
	FUNC_CONCAT_TABLE    = "CONCAT_TABLE"
	FUNC_RETAIN_COLS     = "RETAIN_COLS"
)

type Evaluator struct {
	engine *Engine
	cumSum map[string]float64
}

func NewEvaluator(engine *Engine) *Evaluator {
	return &Evaluator{engine: engine, cumSum: make(map[string]float64)}
}

type EvalFunc func(evaluator *Evaluator, params ...any) (any, error)

var ALL_FUNCTIONS = MergeMaps(MATH_OPERATOR_FUNCTIONS, MATH_FUNCTIONS, SET_FUNCTIONS, STRING_FUNCTIONS, CONTEXT_FUNCTIONS)

var MATH_OPERATOR_FUNCTIONS = map[string]EvalFunc{
	FUNC_PLUS:   _PlusFunc,
	FUNC_MINUS:  _MinusFunc,
	FUNC_MULT:   _MultFunc,
	FUNC_DIVIDE: _DivideFunc,
	FUNC_POW:    _PowFunc,
	FUNC_GT:     _GTFunc,
	FUNC_GE:     _GEFunc,
	FUNC_LT:     _LTFunc,
	FUNC_LE:     _LEFunc,
	FUNC_EQ:     _EqualFunc,
	FUNC_NE:     _NotEqualFunc,
	FUNC_OR:     _OrFunc,
	FUNC_AND:    _AndFunc,
}

var MATH_FUNCTIONS = map[string]EvalFunc{
	FUNC_CEIL:         _CeilFunc,
	FUNC_FLOOR:        _FloorFunc,
	FUNC_ROUND:        _RoundFunc,
	FUNC_MIN:          _MinFunc,
	FUNC_MAX:          _MaxFunc,
	FUNC_MAXV:         _MaxvFunc,
	FUNC_MINV:         _MinvFunc,
	FUNC_COUNT:        _CountFunc,
	FUNC_COUNT_UNIQUE: _CountUniqueFunc,
}

var SET_FUNCTIONS = map[string]EvalFunc{
	FUNC_SET_IN:       _InFunc,
	FUNC_SET_CONTAINS: _ContainsFunc,
}

var STRING_FUNCTIONS = map[string]EvalFunc{
	FUNC_CONCAT:          _ConcatFunc,
	FUNC_STR_LEN:         _StrLenFunc,
	FUNC_STR_CONTAINS:    _StrContainsFunc,
	FUNC_STR_STARTS_WITH: _StartsWithFunc,
	FUNC_STR_ENDS_WITH:   _EndsWithFunc,
	FUNC_IF:              _IfFunc,
	FUNC_EPOCH_MS:        _EpochMsFunc,
	FUNC_NOW_MS:          _NowMsFunc,
	FUNC_SPLIT:           _SplitFunc,
}

var CONTEXT_FUNCTIONS = map[string]EvalFunc{
	FUNC_ROW:            _RowFunc,
	FUNC_LOOKUP:         _LookupFunc,
	FUNC_FILTER:         _FilterFunc,
	FUNC_SUM:            _SumFunc,
	FUNC_AVERAGE:        _AverageFunc,
	FUNC_FIRST:          _FirstFunc,
	FUNC_LAST:           _LastFunc,
	FUNC_JOIN:           _JoinFunc,
	FUNC_CUMULATIVE_SUM: _CumulativeSumFunc,

	FUNC_SYNC_DB_INC:  _SyncDbIncFunc,
	FUNC_SYNC_DB_FULL: _SyncDbFullFunc,
	FUNC_LOAD_CSV:     _LoadCsvFunc,
	FUNC_CARTESIAN:    _CartesianFunc,
	FUNC_XRANGE:       _XrangeFunc,
	FUNC_CONCAT_TABLE: _ConcatTableFunc,
	FUNC_RETAIN_COLS:  _RetainColsFunc,
}

var MATH_OPERATOR_FUNCTION_KEYS = lo.Keys(MATH_OPERATOR_FUNCTIONS)
var SET_FUNCTION_KEYS = lo.Keys(SET_FUNCTIONS)

func _PlusFunc(evaluator *Evaluator, params ...any) (any, error) {
	_r1 := interfaceToFloat(params[0])
	_r2 := interfaceToFloat(params[1])
	_ret := _r1 + _r2
	return _ret, nil
}
func _MinusFunc(evaluator *Evaluator, params ...any) (any, error) {
	return interfaceToFloat(params[0]) - interfaceToFloat(params[1]), nil
}
func _MultFunc(evaluator *Evaluator, params ...any) (any, error) {
	return interfaceToFloat(params[0]) * interfaceToFloat(params[1]), nil
}
func _DivideFunc(evaluator *Evaluator, params ...any) (any, error) {
	v1 := interfaceToFloat(params[0])
	v2 := interfaceToFloat(params[1])
	if v2 == 0 {
		noop("divide by 0")
	}
	return v1 / v2, nil
}
func _PowFunc(evaluator *Evaluator, params ...any) (any, error) {
	return math.Pow(interfaceToFloat(params[0]), interfaceToFloat(params[1])), nil
}
func _CeilFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return nil, nil
	}
	return math.Ceil(interfaceToFloat(vals[0])), nil
}
func _FloorFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return nil, nil
	}
	return math.Floor(interfaceToFloat(vals[0])), nil
}
func _RoundFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return nil, nil
	}
	val := interfaceToFloat(vals[0])
	nDigits := 0
	if len(vals) >= 2 {
		nDigits = int(interfaceToFloat(vals[1]))
	}
	ret := math.Round(val*math.Pow10(nDigits)) / math.Pow10(nDigits)
	return ret, nil
}

func _MinFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) != 2 {
		return math.NaN(), fmt.Errorf("wrong args to MIN")
	}
	v1 := interfaceToFloat(vals[0])
	v2 := interfaceToFloat(vals[1])
	return math.Min(v1, v2), nil
}
func _MaxFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return math.NaN(), fmt.Errorf("MAX needs atleast")
	}

	if len(vals) == 2 {
		v1 := interfaceToFloat(vals[0])
		v2 := interfaceToFloat(vals[1])
		return math.Max(v1, v2), nil
	}
	v1 := interfaceToFloat(vals[0])
	v2 := interfaceToFloat(vals[1])
	return math.Max(v1, v2), nil
}

func _MaxvFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return 0, nil
	}
	var max float64 = 0
	count := 0
	for _, v := range vals {
		val := interfaceToFloat(v)
		if !math.IsNaN(val) {
			if count == 0 {
				max = val
			} else {
				max = math.Max(max, val)
			}
			count++
		}
	}
	return max, nil
}

func _MinvFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return 0, nil
	}
	var min float64 = 0
	count := 0
	for _, v := range vals {
		val := interfaceToFloat(v)
		if !math.IsNaN(val) {
			if count == 0 {
				min = val
			} else {
				min = math.Min(min, val)
			}
			count++
		}
	}
	return min, nil
}

func _EpochMsFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return -1, nil
	}
	val := interfaceToString(vals[0])
	if val == "" {
		return 0, nil
	}
	t1, err := time.Parse(time.DateTime, val)
	if err != nil {
		t2, err2 := time.Parse(time.DateOnly, interfaceToString(val))
		if err2 != nil {
			return -1, nil
		}
		return t2.UnixMilli(), nil
	}
	return t1.UnixMilli(), nil
}

func _NowMsFunc(evaluator *Evaluator, params ...any) (any, error) {
	return time.Now().UnixMilli(), nil
}

func _SplitFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return nil, nil
	}
	str := interfaceToString(vals[0])
	sep := interfaceToString(vals[1])
	res := lo.Map(strings.Split(str, sep), func(item string, index int) ResultValueType { return item })
	return res, nil
}

func _GTFunc(evaluator *Evaluator, params ...any) (any, error) {
	v1 := interfaceToFloat(params[0])
	v2 := interfaceToFloat(params[1])
	ret := v1 > v2
	return ret, nil
}
func _GEFunc(evaluator *Evaluator, params ...any) (any, error) {
	return interfaceToFloat(params[0]) >= interfaceToFloat(params[1]), nil
}
func _LTFunc(evaluator *Evaluator, params ...any) (any, error) {
	return interfaceToFloat(params[0]) < interfaceToFloat(params[1]), nil
}
func _LEFunc(evaluator *Evaluator, params ...any) (any, error) {
	return interfaceToFloat(params[0]) <= interfaceToFloat(params[1]), nil
}
func _EqualFunc(evaluator *Evaluator, params ...any) (any, error) {
	v1 := params[0]
	v2 := convertToType(v1, params[1])
	return v1 == v2, nil
}
func _NotEqualFunc(evaluator *Evaluator, params ...any) (any, error) {
	v1 := params[0]
	v2 := convertToType(v1, params[1])
	return v1 != v2, nil
}
func _OrFunc(evaluator *Evaluator, params ...any) (any, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("bad call to OR function")
	}
	return interfaceToBool(params[0]) || interfaceToBool(params[1]), nil
}
func _AndFunc(evaluator *Evaluator, params ...any) (any, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("bad call to AND function")
	}
	return interfaceToBool(params[0]) && interfaceToBool(params[1]), nil
}

func convertToType(v1, v2 any) any {
	switch v1.(type) {
	case float32:
		return float32(interfaceToFloat(v2))
	case float64:
		return float64(interfaceToFloat(v2))
	case int:
		return int(interfaceToFloat(v2))
	case int32:
		return int32(interfaceToFloat(v2))
	case int64:
		return int64(interfaceToFloat(v2))
	case string:
		return interfaceToString(v2)
	case bool:
		return interfaceToBool(v2)
	default:
		return nil
	}
}

func _InFunc(evaluator *Evaluator, params ...any) (any, error) {
	return nil, nil
}
func _ContainsFunc(evaluator *Evaluator, params ...any) (any, error) {
	return nil, nil
}
func _ConcatFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	valsStr := lo.Map(vals, func(item2 ResultValueType, index int) string { return interfaceToString(item2) })
	_ret := strings.Join(valsStr, "")
	return _ret, nil
}
func _StrLenFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return 0, nil
	}
	v := interfaceToString(vals[0])
	ret := len(v)
	return ret, nil
}
func _StrContainsFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	v1 := interfaceToString(vals[0])
	v2 := interfaceToString(vals[1])
	return strings.Contains(v1, v2), nil
}

func _StartsWithFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	v1 := interfaceToString(vals[0])
	v2 := interfaceToString(vals[1])
	return strings.HasPrefix(v1, v2), nil
}

func _EndsWithFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	v1 := interfaceToString(vals[0])
	v2 := interfaceToString(vals[1])
	return strings.HasSuffix(v1, v2), nil
}

func _IfFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	if val, _ok := vals[0].(bool); _ok && val {
		return vals[1], nil
	}
	return vals[2], nil
}
func _RowFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("Not yet implemented")
}

func _LookupFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("shouldnt be called")
}
func _FilterFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("shouldnt be called")
}

func _SumFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return 0, nil
	}
	var sum float64 = 0
	for _, item := range vals {
		v := interfaceToFloat(item)
		if !math.IsNaN(v) {
			sum += v
		}
	}
	return sum, nil
}

func _AverageFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return 0, nil
	}
	var sum float64 = 0
	count := 0
	for _, item := range vals {
		v := interfaceToFloat(item)
		if !math.IsNaN(v) {
			sum += v
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

func _FirstFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	return vals[0], nil
}

func _LastFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", nil
	}
	return vals[len(vals)-1], nil
}

func _JoinFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) == 0 {
		return "", fmt.Errorf("JOIN needs atleast 2 parameters")
	}
	sep, ok1 := vals[0].(string)
	arr, ok2 := vals[1].([]ResultValueType)
	if !ok1 || !ok2 {
		return "", fmt.Errorf("usage: JOIN <sep> <array>")
	}
	var sb strings.Builder
	for idx, elem := range arr {
		if idx > 0 {
			sb.Write([]byte(sep))
		}
		v := interfaceToString(elem)
		sb.Write([]byte(v))
	}
	ret := sb.String()
	return ret, nil
}

func _CumulativeSumFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	if len(vals) != 2 {
		return "", fmt.Errorf("CUMULATIVE_SUM requires exactly 2 params, got: %v", vals)
	}

	v := interfaceToFloat(vals[0])
	key := interfaceToString(vals[1])
	var newVal float64
	if currVal, found := evaluator.cumSum[key]; !found {
		newVal = v
	} else {
		newVal = v + currVal
	}
	evaluator.cumSum[key] = newVal
	return newVal, nil
}

func _CountFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	return len(vals), nil
}

func _CountUniqueFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	uniq := lo.Uniq(vals)
	return len(uniq), nil
}

func _unrollParams(params ...any) []ResultValueType {
	if len(params) == 0 {
		return []ResultValueType{}
	}
	switch tt := params[0].(type) {
	case []ResultValueType:
		if len(tt) == 0 {
			return []ResultValueType{}
		}
		tt2 := _flatten(tt)
		return tt2
	default:
		return []ResultValueType{}
	}
}

func _flatten(tt []ResultValueType) []ResultValueType {
	if len(tt) == 0 {
		return tt
	}
	switch tt2 := tt[0].(type) {
	case []ResultValueType:
		return tt2
	default:
		return tt
	}
}

func _SyncDbIncFunc(evaluator *Evaluator, params ...any) (any, error) {
	// TODO: Implement
	return "", nil
}

func _SyncDbFullFunc(evaluator *Evaluator, params ...any) (any, error) {
	// TODO: Implement
	return "", nil
}

func _LoadCsvFunc(evaluator *Evaluator, params ...any) (any, error) {
	// TODO: Implement
	return "", nil
}

func _CartesianFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("should not be called")
}

func _XrangeFunc(evaluator *Evaluator, params ...any) (any, error) {
	vals := _unrollParams(params...)
	noop(vals)
	if len(vals) != 3 {
		return nil, fmt.Errorf("bad xrange syntax")
	}
	v1, ok1 := vals[0].(string)
	v2 := interfaceToFloat(vals[1])
	v3 := interfaceToFloat(vals[2])
	if !ok1 {
		return nil, fmt.Errorf("1st param to xrange should be string")
	}
	if math.IsNaN(v2) {
		return nil, fmt.Errorf("2nd param to xrange should be integer")
	}
	if math.IsNaN(v3) {
		return nil, fmt.Errorf("3rd param to xrange should be integer")
	}
	if v2 > v3 {
		return nil, fmt.Errorf("2nd param to xrange should be <= 3rd param")
	}
	nums := make([]ResultValueType, 0, int(v3-v2+1))
	for i := int(v2); i <= int(v3); i++ {
		nums = append(nums, ResultValueType(i))
	}
	table := Table{colNames: []string{v1}, cols: [][]ResultValueType{nums}}
	table.nRows = len(nums)
	return table, nil
}

func _ConcatTableFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("not yet implemented")
}

func _RetainColsFunc(evaluator *Evaluator, params ...any) (any, error) {
	panic("not yet implemented")
}
