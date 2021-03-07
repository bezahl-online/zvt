package instr

import "bezahl.online/zvt/src/apdu/bmp/blen"

type CtrlField struct {
	Class  byte
	Instr  byte
	Length blen.Length
}

// Marshal is
func (c *CtrlField) Marshal(dataLength uint16) []byte {
	var b []byte = []byte{c.Class, c.Instr}
	c.Length.Value = dataLength
	b = append(b, c.Length.Format()...)
	return b
}

// Map maps all used Commands
var Map map[string]CtrlField = make(map[string]CtrlField)

func init() {
	x := CtrlField{
		Class: byte(0x06),
		Instr: byte(0x00),
		Length: blen.Length{
			Kind: byte(blen.BINARY),
		},
	}
	x.Marshal(uint16(1))
	Map["Registration"] = x
	Map["Authorisation"] = CtrlField{
		Class: byte(0x06),
		Instr: byte(0x01),
		Length: blen.Length{
			Kind: byte(blen.BINARY),
		},
	}

}
