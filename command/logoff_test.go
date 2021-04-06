package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestLogOff(t *testing.T) {
	skipShort(t)
	i := instr.Map["ACK"]
	want := Command{
		CtrlField: i,
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	err := PaymentTerminal.LogOff()
	if assert.NoError(t, err) {
		got, err := PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, *got)
		}
	}
}
