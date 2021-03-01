package zvt

import (
	"testing"

	"bezahl.online/zvt/src/zvt/config"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	configByte := config.PaymentReceiptPrintedByECR +
		config.AdminReceiptPrintedByECR +
		config.PTSendsIntermediateStatus +
		config.ECRusingPrintLinesForPrintout
	serviceByte := ServiceMenuNOTAssignedToFunctionKey +
		DisplayTextsForCommandsAuthorisation

	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlv *TLV = &TLV{
		Objects: []DataObject{},
	}
	tlv.Objects = append(tlv.Objects, *listOfCommands)
	want := Response{
		CCRC:   0x84,
		APRC:   0x1E,
		Length: 0x03,
		Data:   []byte{0x6F, 0x09, 0x78},
	}
	got, err := ZVT.Register(&PTConfig{
		config:  byte(configByte),
		service: byte(serviceByte),
		tlv:     tlv,
	})
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}

}

func TestDisplayText(t *testing.T) {
	want := Response{
		CCRC:   0x80,
		APRC:   0x00,
		Length: 0x00,
		Data:   []byte{},
	}
	got, err := ZVT.DisplayText([]string{
		"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
