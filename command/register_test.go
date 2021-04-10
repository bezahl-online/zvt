package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/config"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	skipShort(t) // start
	configByte := config.PaymentReceiptPrintedByECR +
		config.AdminReceiptPrintedByECR +
		config.PTSendsIntermediateStatus +
		// config.AmountInputOnPTNotPossible +
		config.AdminFunctionOnPTNotPossible
	serviceByte := config.Service_MenuNOTAssignedToFunctionKey +
		config.Service_DisplayTextsForCommandsAuthorisationInCAPITALS
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
	i.Length.Size = 1
	err := PaymentTerminal.Register(&Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	})
	assert.NoError(t, err)
}
