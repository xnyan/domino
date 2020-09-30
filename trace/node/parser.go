package node

import (
	//"bytes"
	//"encoding/gob"
	"encoding/binary"
)

func LatInfoToByte(l *L) []byte {
	b := make([]byte, 24, 24)
	binary.LittleEndian.PutUint64(b, uint64(l.S))
	binary.LittleEndian.PutUint64(b[8:], uint64(l.E))
	binary.LittleEndian.PutUint64(b[16:], uint64(l.C))

	return b
}

func ByteToLatInfo(b []byte) *L {
	if len(b) < 24 {
		logger.Fatalf("Invalid []byte insufficient length, expected 24 but %d", len(b))
	}
	l := &L{}
	l.S = int64(binary.LittleEndian.Uint64(b[0:8]))
	l.E = int64(binary.LittleEndian.Uint64(b[8:16]))
	l.C = int64(binary.LittleEndian.Uint64(b[16:24]))

	return l
}

/*
func ToByte(a interface{}) []byte {
	var bf bytes.Buffer
	enc := gob.NewEncoder(&bf)
	enc.Encode(a)
	return bf.Bytes()
}

func ToLatInfo(b []byte) *L {
	bf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bf)
	l := &L{}
	dec.Decode(l)
	return l
}
*/
