package utils

import "encoding/binary"

func BToU16(b []byte, index int, offset int) uint16 {
	return binary.BigEndian.Uint16(b[index:offset])
}

func BToU32(b []byte, index int, offset int) uint32 {
	return binary.BigEndian.Uint32(b[index:offset])
}

func BToU64(b []byte, index int, offset int) uint64 {
	return binary.BigEndian.Uint64(b[index:offset])
}

func U16ToB(v uint16, b []byte) {
	binary.BigEndian.PutUint16(b, v)
}

func U32ToB(v uint32, b []byte) {
	binary.BigEndian.PutUint32(b, v)
}

func U64ToB(v uint64, b []byte) {
	binary.BigEndian.PutUint64(b, v)
}
