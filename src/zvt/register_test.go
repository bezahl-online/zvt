package zvt

import (
	"testing"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/apdu/tlv"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/config"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// start
	configByte := config.PaymentReceiptPrintedByECR +
		config.AdminReceiptPrintedByECR +
		config.PTSendsIntermediateStatus +
		config.ECRusingPrintLinesForPrintout
	serviceByte := config.ServiceMenuNOTAssignedToFunctionKey +
		config.DisplayTextsForCommandsAuthorisation
	var msgSquID *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x1F, 0x73},
		Data: []byte{0, 0, 0},
	}

	var listOfCommands *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlvContainer *tlv.Container = &tlv.Container{
		Objects: []tlv.DataObject{},
	}
	tlvContainer.Objects = append(tlvContainer.Objects, *listOfCommands, *msgSquID)
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
	err := PaymentTerminal.Register(&Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	})
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
		// completion
		i := instr.Map["Completion"]
		want := Command{
			CtrlField: i,
			Data: apdu.DataUnit{
				Data: []byte{},
				BMPOBJs: []bmp.OBJ{
					{ID: 0x19, Data: []byte{0}, Size: 2},
					{ID: 0x29, Data: []byte{0x29, 0x00, 0x10, 0x06}, Size: 5},
					{ID: 0x49, Data: []byte{0x09, 0x78}, Size: 3},
				},
				TLVContainer: tlv.Container{
					Objects: []tlv.DataObject{},
				},
			},
		}
		got, err = PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			if assert.EqualValues(t, want, *got) {
				PaymentTerminal.SendACK()
			}
		}

	}

}
