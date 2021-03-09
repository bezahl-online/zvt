package zvt

import (
	"fmt"
	"testing"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/tlv"
	"github.com/stretchr/testify/assert"
)

func TestAuthorisation(t *testing.T) {
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
	config := &AuthConfig{
		Amount: 1,
	}
	err := ZVT.Authorisation(config)
	got, err := ZVT.ReadResponse()
	if assert.NoError(t, err) {
		if assert.Equal(t, want, *got) {
			for {
				got, err = ZVT.ReadResponse()
				if assert.NoError(t, err) {
					if got.CtrlField.Class == 0x06 && got.CtrlField.Instr == 0x1E {
						if got.Data.Data[0] == 0x6C {
							fmt.Println("Transaction aborted")
							// ZVT.SendACK()
							break
						}
					}
					if got.CtrlField.Class == 0x06 && got.CtrlField.Instr == 0x0F {
						fmt.Println("Transaction successfull")
						ZVT.SendACK()
						break
					}
					ZVT.SendACK()
				}
			}
		}
	}

}
