package blen

import (
	"encoding/binary"
	"fmt"
)

// different kinds of length representations
const (
	// NONE no legth (tag has no data)
	NONE = iota
	// BINARY legth binary coded
	BINARY
	// LL length 0xFx,0xFy -> BCD coded (10x+y)
	LL
	// LLL length 0xFx,0xFy,0xFz -> BCD coded (100x+10y+z)
	LLL
	// BCD fixed length depending on TAG BCD coded
	BCD
)

// Length is the length
type Length struct {
	Kind  byte   // kind of length field
	Size  byte   // size of length field
	Value uint16 // actual data length
}

// Marshal returns a byte representation of the specific length type
func (l *Length) Marshal() ([]byte, error) {
	switch l.Kind {
	case NONE:
		return []byte{}, nil
	case BINARY:
		if l.Value > 254 {
			var b []byte = []byte{0xFF, 0, 0}
			binary.LittleEndian.PutUint16(b[1:], l.Value)

			return b, nil
		}
		return []byte{byte(l.Value)}, nil
	case LL:
		return []byte{
			0xf0 | byte(l.Value/10),
			0xf0 | byte(l.Value-uint16(l.Value/10)*10),
		}, nil
	case LLL:
		z1 := l.Value / 100
		z2 := uint16(uint16(l.Value-z1*100) / 10)
		z3 := uint16(uint16(l.Value - z1*100 - z2*10))
		return []byte{
			0xf0 | byte(z1),
			0xf0 | byte(z2),
			0xf0 | byte(z3),
		}, nil
	}
	return []byte{}, fmt.Errorf("kind of length not implemented")
}

// Unmarshal a length depending on its kind
// fills a Length structure
func (l *Length) Unmarshal(d []byte) error {
	switch l.Kind {
	case BINARY:
		if d[0] == 0xff {
			l.Size = 3
			l.Value = binary.LittleEndian.Uint16(d[1:3])
		} else {
			l.Size = 1
			l.Value = uint16(d[0])
		}
	case LLL:
		if d[0]&0xF0 != 0xF0 ||
			d[1]&0xF0 != 0xF0 ||
			d[2]&0xF0 != 0xF0 {
			return fmt.Errorf("LLL length data invalid: '%03X'", d[:3])
		}
		l.Size = 3
		h := (d[0] & 0x0F)
		z := (d[1] & 0x0F)
		e := (d[2] & 0x0F)
		l.Value = uint16(100*h + 10*z + e)
	case LL:
		if d[0]&0xF0 != 0xF0 ||
			d[1]&0xF0 != 0xF0 {
			return fmt.Errorf("LL length data invalid: '%02X'", d[:2])
		}
		l.Size = 2
		z := (d[0] & 0x0F)
		e := (d[1] & 0x0F)
		l.Value = uint16(10*z + e)
	}
	return nil

}
