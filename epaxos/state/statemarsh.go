package state

import (
	"encoding/binary"
	"io"
)

func (t *Command) Marshal(w io.Writer) {
	var b [8]byte
	bs := b[:8]
	bs = b[:1]
	b[0] = byte(t.Op)
	w.Write(bs)
	//bs = b[:8]
	//binary.LittleEndian.PutUint64(bs, uint64(t.K))
	//w.Write(bs)
	//binary.LittleEndian.PutUint64(bs, uint64(t.V))
	//w.Write(bs)

	marshalStr(string(t.K), w) // key
	marshalStr(string(t.V), w) // value
}

func (t *Command) Unmarshal(r io.Reader) error {
	var b [8]byte
	bs := b[:8]
	bs = b[:1]
	if _, err := io.ReadFull(r, bs); err != nil {
		return err
	}
	t.Op = Operation(b[0])
	//bs = b[:8]
	//if _, err := io.ReadFull(r, bs); err != nil {
	//	return err
	//}
	//t.K = Key(binary.LittleEndian.Uint64(bs))
	//if _, err := io.ReadFull(r, bs); err != nil {
	//	return err
	//}
	//t.V = Value(binary.LittleEndian.Uint64(bs))

	// Key
	if s, err := unmarshalStr(r); err != nil {
		return err
	} else {
		t.K = Key(s)
	}

	// Value
	if s, err := unmarshalStr(r); err != nil {
		return err
	} else {
		t.V = Value(s)
	}

	return nil
}

func (t *Key) Marshal(w io.Writer) {
	//var b [8]byte
	//bs := b[:8]
	//binary.LittleEndian.PutUint64(bs, uint64(*t))
	//w.Write(bs)
	marshalStr(string(*t), w)
}

func (t *Value) Marshal(w io.Writer) {
	//var b [8]byte
	//bs := b[:8]
	//binary.LittleEndian.PutUint64(bs, uint64(*t))
	//w.Write(bs)
	marshalStr(string(*t), w)
}

func (t *Key) Unmarshal(r io.Reader) error {
	//var b [8]byte
	//bs := b[:8]
	//if _, err := io.ReadFull(r, bs); err != nil {
	//	return err
	//}
	//*t = Key(binary.LittleEndian.Uint64(bs))
	if s, err := unmarshalStr(r); err != nil {
		return err
	} else {
		*t = Key(s)
	}
	return nil
}

func (t *Value) Unmarshal(r io.Reader) error {
	//var b [8]byte
	//bs := b[:8]
	//if _, err := io.ReadFull(r, bs); err != nil {
	//	return err
	//}
	//*t = Value(binary.LittleEndian.Uint64(bs))
	if s, err := unmarshalStr(r); err != nil {
		return err
	} else {
		*t = Value(s)
	}
	return nil
}

func marshalStr(s string, w io.Writer) {
	bArr := []byte(s)
	len := uint64(len(bArr))
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, len)
	w.Write(bs)
	if len > 0 {
		w.Write(bArr)
	}
}

func unmarshalStr(r io.Reader) (string, error) {
	var b [8]byte
	bs := b[:8]
	if _, err := io.ReadFull(r, bs); err != nil {
		return "", err
	}
	len := binary.LittleEndian.Uint64(bs)
	if len > 0 {
		sb := make([]byte, len, len)
		bs = sb[:len]
		if _, err := io.ReadFull(r, bs); err != nil {
			return "", err
		}
		return string(bs), nil
	}
	return "", nil
}
