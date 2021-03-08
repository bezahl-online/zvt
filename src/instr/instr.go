package instr

import "bezahl.online/zvt/src/apdu/bmp/blen"

// CtrlField is
type CtrlField struct {
	Class         byte
	Instr         byte
	Length        blen.Length
	RawDataLength int
}

// Marshal is
func (c *CtrlField) Marshal(dataLength uint16) []byte {
	var b []byte = []byte{c.Class, c.Instr}
	c.Length.Value = dataLength
	b = append(b, c.Length.Format()...)
	return b
}

// Find searches for Class and Instr in the Map
func Find(command *[]byte) *CtrlField {
	for key := range Map {
		if Map[key].Class == (*command)[0] &&
			Map[key].Instr == (*command)[1] {
			c := Map[key]
			return &c
		}
	}
	return nil
}

// Map maps all used Commands
var Map map[string]CtrlField = make(map[string]CtrlField)

func init() {
	Map["Registration"] = CtrlField{
		Class: byte(0x06),
		Instr: byte(0x00),
		Length: blen.Length{
			Kind: byte(blen.BINARY),
		},
	}
	Map["Authorisation"] = CtrlField{
		Class: byte(0x06),
		Instr: byte(0x01),
		Length: blen.Length{
			Kind: byte(blen.BINARY),
		},
	}
	Map["ACK"] = CtrlField{
		Class: byte(0x80),
		Instr: byte(0x00),
	}
	Map["Completion"] = CtrlField{
		Class: byte(0x06),
		Instr: byte(0x0F),
		Length: blen.Length{
			Kind:  blen.BINARY,
			Value: 0,
		},
	}
	Map["Intermediate"] = CtrlField{
		Class:         byte(0x04),
		Instr:         byte(0xFF),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 2,
	}
	Map["AccountBalance"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x03),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["PrintTextBlock"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0xD3),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["NotSupported"] = CtrlField{
		Class:         byte(0x84),
		Instr:         byte(0x83),
		Length:        blen.Length{Kind: blen.NONE, Value: 0},
		RawDataLength: 0,
	}

}
