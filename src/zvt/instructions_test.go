package zvt

import (
	"testing"
	"time"

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
	var msgSquID *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x73},
		data: []byte{0, 0, 0},
	}

	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlv *TLV = &TLV{
		Objects: []DataObject{},
	}
	tlv.Objects = append(tlv.Objects, *listOfCommands, *msgSquID)
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

func TestAuthorisation(t *testing.T) {
	want := Response{
		CCRC:   0x04,
		APRC:   0xFF,
		Length: 0x2E,
		Data:   []byte{0x0A, 0x01, 0x06, 0x2A},
	}
	// var cardPollTimeout *DataObject = &DataObject{
	// 	TAG:  []byte{0x1F, 0x5B},
	// 	data: []byte{0x10},
	// }
	var msgSquID *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x73},
		data: []byte{0, 0, 0},
	}
	var objects *[]DataObject = &[]DataObject{}

	*objects = append(*objects, *msgSquID)
	// var paymentType byte = payment.PaymentIncludeGeldKarte + payment.PrinterReady + payment.GirocardTransaction
	currency := EUR
	var tlv *TLV = &TLV{
		Objects: *objects,
	}
	config := &AuthConfig{
		Amount:   1,
		Currency: &currency,
		// PaymentType: &paymentType,
		TLV: tlv,
	}
	got, err := ZVT.Authorisation(config)
	got.Data = got.Data[:4]
	if assert.NoError(t, err) {
		if assert.EqualValues(t, want, *got) {
			ZVT.SendACK(5 * time.Second)
		}
	}

}
