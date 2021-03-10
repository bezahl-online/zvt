package zvt

import (
	"testing"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/apdu/tlv"
	"bezahl.online/zvt/src/instr"
	"github.com/stretchr/testify/assert"
)

func TestLogOff(t *testing.T) {
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
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}
