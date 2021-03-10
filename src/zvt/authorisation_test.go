package zvt

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bezahl-online/zvt/src/apdu"
	"github.com/bezahl-online/zvt/src/apdu/bmp"
	"github.com/bezahl-online/zvt/src/apdu/tlv"
	"github.com/bezahl-online/zvt/src/instr"
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
	err := PaymentTerminal.Authorisation(config)
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		if assert.Equal(t, want, *got) {
			done := false
			for !done {
				got, err = PaymentTerminal.ReadResponseWithTimeout(20 * time.Second)
				if assert.NoError(t, err) {
					PaymentTerminal.SendACK()
					switch got.CtrlField.Class {
					case 0x06:
						switch got.CtrlField.Instr {
						case 0x1E:
							if got.Data.Data[0] == 0x6C {
								fmt.Println("Transaction aborted")
								done = true
							}
						case 0x0F:
							fmt.Println("Transaction successfull")
							done = true
						}
					}
				} else {
					log.Fatal(err)
				}
			}
		}
	}

}
