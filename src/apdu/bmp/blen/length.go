package blen

import "encoding/binary"

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

type Length struct {
	Kind  byte
	Value uint16
}

// Format returns a byte representation of the specific length type
func (l *Length) Format() []byte {
	switch l.Kind {
	case NONE:
		return []byte{}
	case BINARY:
		if l.Value > 254 {
			var b []byte = []byte{0, 0}
			binary.LittleEndian.PutUint16(b, l.Value)
			return b
		}
		return []byte{byte(l.Value)}
	case LL:
		return []byte{
			0xf0 | byte(l.Value/10),
			0xf0 | byte(l.Value-uint16(l.Value/10)*10),
		}
	case LLL:
		z1 := l.Value / 100
		z2 := uint16(uint16(l.Value-z1*100) / 10)
		z3 := uint16(uint16(l.Value - z1*100 - z2*10))
		return []byte{
			0xf0 | byte(z1),
			0xf0 | byte(z2),
			0xf0 | byte(z3),
		}
	}
	return []byte{}
}
