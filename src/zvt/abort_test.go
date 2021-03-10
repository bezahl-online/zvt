package zvt

import (
	"testing"

	"github.com/bezahl-online/zvt/src/apdu"
	"github.com/bezahl-online/zvt/src/apdu/bmp"
	"github.com/bezahl-online/zvt/src/apdu/tlv"
	"github.com/bezahl-online/zvt/src/instr"
	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
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
	err := PaymentTerminal.Abort()
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}
