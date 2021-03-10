package zvt

import (
	"fmt"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/tlv"
)

// Command is the structur for a APDU
type Command struct {
	CtrlField instr.CtrlField
	Data      apdu.DataUnit
}

// Marshal marshal every thing to final command
func (c *Command) Marshal() ([]byte, error) {
	var b []byte = []byte{}
	data, err := c.Data.Marshal()
	if err != nil {
		return b, err
	}
	cf, err := c.CtrlField.Marshal(uint16(len(data)))
	if err != nil {
		return b, err
	}
	b = append(b, cf...)
	b = append(b, data...)
	return b, nil
}

// Unmarshal is
func (c *Command) Unmarshal(data *[]byte) error {
	i := instr.Find(data)
	if i == nil {
		return fmt.Errorf("APRC %0X not found", (*data)[:2])
	}
	c.CtrlField = *i
	dstart := 3
	if (*data)[2] == 0xff {
		dstart = 5
	}
	dend := dstart + i.RawDataLength
	if (*data)[2] < byte(i.RawDataLength) {
		dend = dstart + int((*data)[2])
	}
	raw := (*data)[dstart:dend]
	objs := []bmp.OBJ{}
	x := 100
	for ; x > 0; x-- {
		if len(*data) <= dend || (*data)[dend] == tlv.BMPTLV {
			break
		}
		o := bmp.OBJ{}
		err := o.Unmarshal(((*data)[dend:]))
		if err != nil {
			return err
		}
		objs = append(objs, o)
		dend += o.Size
	}
	if x < 1 {
		return fmt.Errorf("loop exceeded 100 iterations while unmashalling data objects")
	}
	tlv := tlv.Container{
		Objects: []tlv.DataObject{},
	}
	if len(*data) < dend {
		return nil
	}
	tlvData := (*data)[dend:]
	err := tlv.Unmarshal(&tlvData)
	if err != nil {
		return err
	}
	c.Data = apdu.DataUnit{
		Data:         raw,
		BMPOBJs:      objs,
		TLVContainer: tlv,
	}
	return nil
}
