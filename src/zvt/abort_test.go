package zvt

import (
	"testing"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/tlv"
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
	err := ZVT.Abort()
	got, err := ZVT.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}
