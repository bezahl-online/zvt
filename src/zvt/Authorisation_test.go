package zvt

// import (
// 	"testing"
// 	"time"

// 	"bezahl.online/zvt/src/zvt/payment"
// 	"bezahl.online/zvt/src/zvt/tlv"
// 	"github.com/stretchr/testify/assert"
// )

// func TestAuthorisation(t *testing.T) {
// 	// var cardPollTimeout *tlv.DataObject = &tlv.DataObject{
// 	// 	TAG:  []byte{0x1F, 0x5B},
// 	// 	data: []byte{0x10},
// 	// }
// 	return // FIXME: total überarbeiten!
// 	var msgSquID *tlv.DataObject = &tlv.DataObject{
// 		TAG:  []byte{0x1F, 0x73},
// 		Data: []byte{0, 0, 0},
// 	}
// 	var objects *[]tlv.DataObject = &[]tlv.DataObject{}

// 	*objects = append(*objects, *msgSquID)
// 	var paymentType byte = payment.PrinterReady + payment.GirocardTransaction
// 	currency := EUR
// 	var tlv *tlv.Container = &tlv.Container{
// 		Objects: *objects,
// 	}
// 	config := &AuthConfig{
// 		Amount:      1,
// 		Currency:    &currency,
// 		PaymentType: &paymentType,
// 		TLV:         tlv,
// 	}
// 	got, err := ZVT.Authorisation(config)
// 	got.Data = got.Data[:4]
// 	if assert.NoError(t, err) {
// 		if assert.Equal(t, true, got.IsACK()) {
// 			for {
// 				got, err = ZVT.readResponse(5 * time.Second)
// 				if assert.NoError(t, err) {
// 					// if assert.Equal(t, true, got.IsStatus() || got.IsIntermediate()) {
// 					ZVT.SendACK(5 * time.Second)
// 					// }
// 				}
// 			}
// 		}
// 	}

// }
