package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
)

// Good read about Arrow: https://www.apachecon.com/acna2022/slides/01_Topol_Arrow_and_Go.pdf

// func Test_Arrow_StringArray(t *testing.T) {
// 	pool := memory.NewGoAllocator()

// 	lb := array.NewStringBuilder(pool)
// 	defer lb.Release()

// 	lb.Append("str1")
// 	lb.Append("str2")
// 	lb.Append("str3")
// 	lb.Append("str4")
// 	lb.Append("str5")

// 	arr := lb.NewArray().(*array.String)
// 	defer arr.Release()

// 	fmt.Printf("NullN()   = %d\n", arr.NullN())
// 	fmt.Printf("Len()     = %d\n", arr.Len())
// 	fmt.Printf("Type()    = %v\n", arr.DataType())
// 	fmt.Printf("List      = %v\n", arr)

// 	noop(pool)
// }

func _makeArrays() ([]int, [][]string) {
	// nValsArr := []int{100 * 1000, 200 * 1000, 500 * 1000, 1000 * 1000}
	nValsArr := []int{1000 * 1000}
	arr := make([][]string, 0, len(nValsArr))
	for _, numVals := range nValsArr {
		vals := lo.Map(Xrange(0, numVals), func(item int, index int) string { return fmt.Sprintf("str_%d", item) })
		arr = append(arr, vals)
	}
	return nValsArr, arr
}

func Test_Serialize_Arrow_StringArray(t *testing.T) {
	nValsArr, arr := _makeArrays()

	for idx, numVals := range nValsArr {
		vals := arr[idx]

		ts := time.Now()
		bytes := SerializeArrowFormat(vals)
		fmt.Println("num bytes: ", numVals, len(bytes))
		fmt.Println("time taken in serialization: ", time.Since(ts).Milliseconds())
	}

	bytes := SerializeArrowFormat(arr[len(arr)-1])

	ts := time.Now()
	ret, err := DeserializeArrowFormat(bytes)
	if err != nil {
		t.Fatalf("error deserializing: %s", err)
	}

	noop(ret)
	fmt.Println("time taken in deser: ", time.Since(ts).Milliseconds())
}

func Test_Serialize_SimpleStringArray(t *testing.T) {
	nValsArr, arr := _makeArrays()
	for idx, numVals := range nValsArr {
		vals := arr[idx]

		ts := time.Now()
		bytes := _serializeSimpleStringArr(vals)
		fmt.Println("num bytes: ", numVals, len(bytes))
		fmt.Println("time taken: ", time.Since(ts).Milliseconds())
	}
}

func _serializeSimpleStringArr(vals []string) []byte {
	bb := bytes.Buffer{}
	for _, v := range vals {
		bb.WriteString(v)
	}
	for _, v := range vals {
		binary.Write(&bb, binary.LittleEndian, int32(len(v)))
	}

	bytes := bb.Bytes()
	return bytes
}
