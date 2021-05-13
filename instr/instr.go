package instr

import "github.com/bezahl-online/zvt/apdu/bmp/blen"

// CtrlField is
type CtrlField struct {
	Class         byte
	Instr         byte
	Length        blen.Length
	RawDataLength int
}

// Marshal is
func (c *CtrlField) Marshal(dataLength uint16) ([]byte, error) {
	var b []byte = []byte{c.Class, c.Instr}
	c.Length.Value = dataLength
	l, err := c.Length.Marshal()
	if err != nil {
		return b, err
	}
	b = append(b, l...)
	return b, nil
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
	Map["Status"] = CtrlField{
		Class:         byte(0x05),
		Instr:         byte(0x01),
		Length:        blen.Length{Kind: byte(blen.BINARY)},
		RawDataLength: 0,
	}
	Map["Registration"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x00),
		Length:        blen.Length{Kind: byte(blen.BINARY)},
		RawDataLength: 0,
	}
	Map["Authorisation"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x01),
		Length:        blen.Length{Kind: byte(blen.BINARY)},
		RawDataLength: 0,
	}
	Map["LogOff"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x02),
		Length:        blen.Length{Kind: byte(blen.BINARY), Value: 0},
		RawDataLength: 0,
	}
	Map["AccountBalance"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x03),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["TransactionAborted"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x1E),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 1, // password
	}
	Map["EndOfDay"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x50),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["PrintLine"] = CtrlField{ // NOT realy implemented correctly
		Class:         byte(0x06),
		Instr:         byte(0xD1),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 1, // FIXME: actually there is only raw data
		// <attribut><text>
	}
	Map["Abort"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0xB0),
		Length:        blen.Length{Kind: byte(blen.BINARY), Value: 0},
		RawDataLength: 0,
	}
	Map["PrintTextBlock"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0xD3),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["DisplayText"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0xE0),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 1,
	}
	Map["ACK"] = CtrlField{
		Class:         byte(0x80),
		Instr:         byte(0x00),
		Length:        blen.Length{Kind: blen.BINARY, Size: 0, Value: 0},
		RawDataLength: 0,
	}
	Map["NAK"] = CtrlField{
		Class:         byte(0x84),
		Instr:         byte(0x9C),
		Length:        blen.Length{Kind: blen.BINARY, Size: 0, Value: 0},
		RawDataLength: 0,
	}

	// mostly from PT
	Map["Intermediate"] = CtrlField{
		Class:         byte(0x04),
		Instr:         byte(0xFF),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 2,
	}
	Map["StatusInformation"] = CtrlField{
		Class:         byte(0x04),
		Instr:         byte(0x0F),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["Completion"] = CtrlField{
		Class:         byte(0x06),
		Instr:         byte(0x0F),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["AbortAuthorisation"] = CtrlField{
		Class:         byte(0x84),
		Instr:         byte(0xA0),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}
	Map["NotSupported"] = CtrlField{
		Class:         byte(0x84),
		Instr:         byte(0x83),
		Length:        blen.Length{Kind: blen.BINARY, Value: 0},
		RawDataLength: 0,
	}

}
