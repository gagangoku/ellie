package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/samber/lo"
)

type LoaderInf interface {
	GetColumn(col string) (ColMetdata, []ResultValueType, bool)
	SaveColumn(metadata ColMetdata, vals []ResultValueType, save bool) error
}

type FileLoader struct {
	backupDir string
	appId     string
}

func (l *FileLoader) Init(backupDir, appId string) {
	l.backupDir = backupDir
	l.appId = appId
}

func (l *FileLoader) GetColumn(col string) (ColMetdata, []ResultValueType, bool) {
	if l.backupDir == "" {
		return ColMetdata{}, nil, false
	}

	filename := fmt.Sprintf("%s/%s/%s", l.backupDir, l.appId, col)
	bytes, err := os.ReadFile(filename)
	if err != nil || len(bytes) < 2 {
		return ColMetdata{}, nil, false
	}

	nb := int(binary.BigEndian.Uint16(bytes[0:2]))
	if nb == 0 || len(bytes) < nb+1 {
		return ColMetdata{}, nil, false
	}
	mObj := ColMetdata{}
	jsonStr := bytes[2 : nb+2]
	err = json.Unmarshal(jsonStr, &mObj)
	if err != nil {
		return ColMetdata{}, nil, false
	}

	colBytes := bytes[nb+2:]
	vals, err := DeserializeArrowFormat(colBytes)
	if err != nil {
		return ColMetdata{}, nil, false
	}

	v := lo.Map(vals, func(item string, index int) ResultValueType { return item })
	return mObj, v, true
}

func (l *FileLoader) SaveColumn(metadata ColMetdata, vals []ResultValueType, save bool) error {
	if l.backupDir == "" || !save {
		return nil
	}

	_bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(len(_bytes)))

	buf.Write(b)
	buf.Write(_bytes)

	valStrings := lo.Map(vals, func(item ResultValueType, index int) string { return interfaceToString(item) })
	_bytes = SerializeArrowFormat(valStrings)
	buf.Write(_bytes)

	fBytes := buf.Bytes()

	// Create app backup dir if it doesnt exist
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", l.backupDir, l.appId), os.ModePerm); err != nil {
		return fmt.Errorf("couldnt create backup dir: %s", err)
	}

	filename := fmt.Sprintf("%s/%s/%s", l.backupDir, l.appId, metadata.Column)
	err = os.WriteFile(filename, fBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %s %s", filename, err)
	}
	return nil
}

func DeserializeArrowFormat(byteArr []byte) ([]string, error) {
	pool := memory.NewGoAllocator()
	rb := array.NewRecordBuilder(pool, _arrowSchema)
	defer rb.Release()

	buf := bufio.NewReader(bytes.NewBuffer(byteArr))
	reader, err := ipc.NewReader(buf, ipc.WithAllocator(pool), ipc.WithSchema(_arrowSchema))
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	for reader.Next() {
		rec := reader.Record()
		vals := rec.Column(0).(*array.String)
		for i := 0; i < vals.Len(); i++ {
			ret = append(ret, vals.Value(i))
		}
	}
	return ret, nil
}

func SerializeArrowFormat(vals []string) []byte {
	pool := memory.NewGoAllocator()
	rb := array.NewRecordBuilder(pool, _arrowSchema)
	defer rb.Release()

	rb.Field(0).(*array.StringBuilder).AppendValues(vals, nil)

	rec := rb.NewRecord()
	defer rec.Release()

	bb := bytes.Buffer{}
	bw := bufio.NewWriter(&bb)
	writer := ipc.NewWriter(bw, ipc.WithAllocator(pool), ipc.WithSchema(_arrowSchema))
	writer.Write(rec)
	bw.Flush()

	bytes := bb.Bytes()
	return bytes
}

var _arrowSchema = arrow.NewSchema(
	[]arrow.Field{
		{Name: "col1", Type: arrow.BinaryTypes.String},
	},
	nil,
)
