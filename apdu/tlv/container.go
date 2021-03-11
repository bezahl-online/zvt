package tlv

import (
	"bytes"
	"fmt"
)

// BMPTLV is 0x06
const BMPTLV = byte(6)

// Container is the type length value container
type Container struct {
	Objects []DataObject
}

// Marshal retuns the byte array of the tlv
func (t *Container) Marshal() []byte {
	var b []byte
	if len(t.Objects) > 0 {
		b = append(b, BMPTLV)
		data := MarshalDataObjects(&t.Objects)
		b = append(b, CompileLength(len(data))...)
		b = append(b, data...)
	}
	return b
}

// Unmarshal fills the structur with the given data
// data must include TLV BMP/TAG 0x06
func (t *Container) Unmarshal(data *[]byte) error {
	d := *data
	idx := bytes.IndexByte(d, 0x06)
	if idx >= 0 && len(d) > 3 {
		d = d[idx:]
		lenData := d[1:5]
		tlvLen, sizeOfLenField, err := DecompileLength(&lenData)
		if err != nil {
			return err
		}
		realLen := uint16(len(d)) - sizeOfLenField - 1
		if realLen != tlvLen {
			return fmt.Errorf("value in length field (%d) and length of data (%d) does not match", tlvLen, realLen)
		}
		// reduce data to data after TLV TAG (06) and length byte(s)
		d = d[sizeOfLenField+1:]
		if t.Objects == nil {
			t.Objects = []DataObject{}
		}
		for {
			var obj DataObject = DataObject{}
			objLength, err := obj.Unmarshal(d)
			if err != nil {
				return err
			}
			t.Objects = append(t.Objects, obj)
			if len(d) == int(objLength) {
				break
			}
			d = d[objLength:]
		}
	}
	return nil
}
