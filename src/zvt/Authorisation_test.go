package zvt

// import (
// 	"testing"
// 	"time"

// 	"bezahl.online/zvt/src/apdu"
// 	"bezahl.online/zvt/src/apdu/bmp"
// 	"bezahl.online/zvt/src/instr"
// 	"bezahl.online/zvt/src/zvt/payment"
// 	"bezahl.online/zvt/src/zvt/tlv"
// 	"github.com/stretchr/testify/assert"
// )

// func TestAuthorisation(t *testing.T) {
// 	// var cardPollTimeout *tlv.DataObject = &tlv.DataObject{
// 	// 	TAG:  []byte{0x1F, 0x5B},
// 	// 	data: []byte{0x10},
// 	// }
// 	var msgSquID *tlv.DataObject = &tlv.DataObject{
// 		TAG:  []byte{0x1F, 0x73},
// 		Data: []byte{0, 0, 0},
// 	}
// 	var objects *[]tlv.DataObject = &[]tlv.DataObject{}

// 	*objects = append(*objects, *msgSquID)
// 	var paymentType byte = payment.PrinterReady + payment.GirocardTransaction
// 	currency := EUR
// 	var tlvContainer *tlv.Container = &tlv.Container{
// 		Objects: *objects,
// 	}
// 	i := instr.Map["ACK"]
// 	want := Command{
// 		Instr: i,
// 		Data: apdu.DataUnit{
// 			Data:         []byte{},
// 			BMPOBJs:      []bmp.OBJ{},
// 			TLVContainer: tlv.Container{},
// 		},
// 	}
// 	config := &AuthConfig{
// 		Amount:      1,
// 		Currency:    &currency,
// 		PaymentType: &paymentType,
// 		TLV:         tlvContainer,
// 	}
// 	err := ZVT.Authorisation(config)
// 	got, err := ZVT.ReadResponse(time.Second * 5)
// 	if assert.NoError(t, err) {
// 		if assert.Equal(t, want, got) {
// 			for {
// 				got, err = ZVT.ReadResponse(5 * time.Second)
// 				if assert.NoError(t, err) {
// 					// if assert.Equal(t, true, got.IsStatus() || got.IsIntermediate()) {
// 					ZVT.SendACK()
// 					// }
// 				}
// 			}
// 		}
// 	}

// }
