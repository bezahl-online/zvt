package command

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/apdu/tlv/tag"
	"github.com/bezahl-online/zvt/instr"
	"go.uber.org/zap"
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
	var err error
	i := instr.Find(data)
	if i == nil {
		return fmt.Errorf("APRC %0X not found", (*data)[:2])
	}
	c.CtrlField = *i
	// FIXME: workaround für end_of_day bug
	if c.CtrlField.Class == 0x04 &&
		c.CtrlField.Instr == 0x0f &&
		(*data)[2] == 0x27 {
		// length field missing bug
		// we will insert it
		l := len(*data)
		if *data, err = insert(data, 2, byte(l-2)); err != nil {
			return err
		}
		Logger.Debug("Workaround for missing length field",
			zap.String("Data", fmt.Sprintf("a length of %d is inserted", l)))
	}
	if c.CtrlField.Class == 0x06 &&
		c.CtrlField.Instr == 0x0F &&
		(*data)[2] == 0x01 {
		if *data, err = insert(data, 3, 0x27); err != nil {
			return err
		}
		(*data)[2] = 0x02
	}
	if err := c.CtrlField.Length.Unmarshal((*data)[2:]); err != nil {
		return err
	}
	// after binary length field
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
		Logger.Debug(fmt.Sprintf("BMP %02X: '%s'", o.ID, bmp.InfoMap[o.ID].Name))
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
	err = tlv.Unmarshal(&tlvData)
	if err != nil {
		return err
	}
	for _, o := range tlv.Objects {
		var tagID [2]byte = [2]byte{o.TAG[0]}
		if tagID[0] == 0x1F {
			tagID[1] = o.TAG[1]
		}
		Logger.Debug(fmt.Sprintf("TAG %02X: '%s'", o.TAG, tag.InfoMap[tagID].Name))
	}
	c.Data = apdu.DataUnit{
		Data:         raw,
		BMPOBJs:      objs,
		TLVContainer: tlv,
	}
	return nil
}

func insert(data *[]byte, pos int, b byte) ([]byte, error) {
	if pos > len(*data) {
		return []byte{}, fmt.Errorf("position greater than length of data")
	}
	var d []byte = make([]byte, len(*data)+1)
	copy(d, *data)
	d = append(d[:pos], b)
	d = append(d, (*data)[pos:]...)
	return d, nil
}

// IsAck returns nil if it is ACK else  error
func (c *Command) IsAck() (err error) {
	isAck := c.CtrlField.Class == 0x80 && c.CtrlField.Instr == 0x00
	if !isAck {
		err = fmt.Errorf("error % 02X %+v", []byte{c.CtrlField.Class, c.CtrlField.Instr}, c.Data)
		Logger.Error(err.Error())
	}
	return err
}
