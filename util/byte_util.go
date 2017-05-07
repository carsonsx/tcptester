package util

import "encoding/binary"

func AddUint8ToBytePrefix(raw []byte, prefix uint8) []byte {
	return append(raw, byte(prefix))
}

func AddUint16ToBytePrefix(raw []byte, prefix uint16, littleEndian bool) []byte {
	l := 2
	data := make([]byte, l + len(raw))
	if littleEndian {
		binary.LittleEndian.PutUint16(data, prefix)
	} else {
		binary.BigEndian.PutUint16(data, prefix)
	}
	copy(data[l:], raw)
	return data
}

func AddUint32ToBytePrefix(raw []byte, prefix uint32, littleEndian bool) []byte {
	l := 4
	data := make([]byte, l + len(raw))
	if littleEndian {
		binary.LittleEndian.PutUint32(data, prefix)
	} else {
		binary.BigEndian.PutUint32(data, prefix)
	}
	copy(data[l:], raw)
	return data
}
