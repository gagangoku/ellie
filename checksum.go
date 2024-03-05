package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

func Checksum(vals []ResultValueType) string {
	buf := bytes.Buffer{}
	for _, v := range vals {
		val := interfaceToString(v)
		buf.Write([]byte(val))
		buf.WriteByte(byte(len(val)))
	}

	// crc32q := crc32.MakeTable(0xD5828281)
	crc32q := crc32.MakeTable(crc32.Castagnoli)
	crc := crc32.Checksum(buf.Bytes(), crc32q)

	// Add the number of vals to checksum. Had a case where checksum matched with different number of rows
	return fmt.Sprintf("%d-%d", len(vals), crc)
}
