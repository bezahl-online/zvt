package command

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
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
	// after binary length field
	dstart := 3
	if (*data)[2] == 0xff {
		dstart = 5
	}
	// FIXME: workaround für end_of_day bug
	if c.CtrlField.Class == 0x04 &&
		c.CtrlField.Instr == 0x0f &&
		(*data)[2] == 0x27 {
		// length field missing bug
		// we will insert it
		var d []byte = make([]byte, len(*data)+1)
		copy(d, *data)
		l := len(*data)
		d = append(d[:2], byte(l-2))
		d = append(d, (*data)[2:]...)
		*data = d
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

// IsAck returns true if command is ack
func (c *Command) IsAck() bool {
	return c.CtrlField.Class == 0x80 && c.CtrlField.Instr == 0x00
}
