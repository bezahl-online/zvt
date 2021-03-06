package zvt

import (
	"testing"

	"bezahl.online/zvt/src/zvt/payment"
	"bezahl.online/zvt/src/zvt/tlv"
	"bezahl.online/zvt/src/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestGetBytes(t *testing.T) {
	c := Command{
		Class: 01,
		Inst:  02,
		Data:  []byte("Teststring"),
	}
	got := c.compile()
	want := []byte{0x01, 0x02, 0x0a, 0x54, 0x65, 0x73,
		0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	assert.Equal(t, want, got)
}

func TestCompileText(t *testing.T) {
	want := []byte{0xf1, 0xf0, 0xf4, 0x54, 0x65, 0x73,
		0x74, 0xf2, 0xf0, 0xf5, 0x41, 0x72, 0x72, 0x61, 0x79}
	got := compileText([]string{"Test", "Array"})
	assert.Equal(t, want, got)
}

func TestCompilePTConfig(t *testing.T) {
	want := []byte{0x12, 0x34, 0x56, 0x8E, 0x9, 0x78, 0x3,
		0x3, 0x6, 0x6, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}
	var listOfCommands *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlvContainer *tlv.Container = &tlv.Container{
		Objects: []tlv.DataObject{},
	}
	tlvContainer.Objects = append(tlvContainer.Objects, *listOfCommands)
	config := Config{
		pwd:          [3]byte{0x12, 0x34, 0x56},
		config:       0x8E,
		currency:     EUR,
		service:      0x03,
		tlvContainer: tlvContainer,
	}
	got := config.CompileConfig()
	assert.EqualValues(t, want, got)
}

func TestUnmarshalAPDU(t *testing.T) {
	testBytes := []byte{0x84, 0x1E, 0x03, 0x6F, 0x09, 0x97}
	want := Response{
		CCRC:   0x84,
		APRC:   0x1E,
		Length: 0x03,
		Data:   []byte{0x6F, 0x09, 0x97},
	}
	got, err := ZVT.unmarshalAPDU(testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}

func TestUnmarshalAPDUData(t *testing.T) {
	t.Run("from data files", func(t *testing.T) {
		testBytes, err := util.Load("dump/data050730027.bin")
		if !assert.NoError(t, err) {
			return
		}
		want := Response{
			CCRC:   0x04,
			APRC:   0xFF,
			Length: 0x2E,
			Data:   testBytes[3:],
		}

		var got *Response
		got, err = ZVT.unmarshalAPDU(testBytes)
		if assert.NoError(t, err) {
			if assert.EqualValues(t, want, *got) {
				want := tlv.Container{
					Objects: []tlv.DataObject{},
				}
				want.Objects = append(want.Objects, tlv.DataObject{
					TAG:  []byte{0x24},
					Data: testBytes[9:],
				})
				got, err := (*got).GetTLV()
				if assert.NoError(t, err) {
					assert.EqualValues(t, want, *got)
				}
			}

		}
	})
	// t.Run("from data files", func(t *testing.T) {
	// 	testBytes, err := util.Load("dump/data051327012.bin")
	// 	if !assert.NoError(t, err) {
	// 		return
	// 	}
	// 	want := Response{
	// 		CCRC:    0x06,
	// 		APRC:    0xD3,
	// 		Length:  0xFF,
	// 		IStatus: 0,
	// 		Data:    testBytes,
	// 		TLV: TLV{
	// 			Objects: []DataObject{},
	// 		},
	// 	}
	// 	if len(testBytes) < 9 {
	// 		fmt.Printf("% X", testBytes)
	// 		return
	// 	}
	// 	want.TLV.Objects = append(want.TLV.Objects, DataObject{
	// 		TAG:  []byte{0x24},
	// 		data: testBytes[9:],
	// 	})
	// 	var got *Response
	// 	got, err = ZVT.unmarshalAPDU(testBytes)
	// 	if assert.NoError(t, err) {
	// 		assert.EqualValues(t, want, *got)
	// 	}
	// })
}

func TestAuthData(t *testing.T) {
	want := []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x49, 0x09,
		0x78, 0x19, 0x75, 0x0E, 0x11, 0x21, 0x06, 0x04, 0x1F, 0x5B, 0x01, 0x05}
	var cardPollTimeout *tlv.DataObject = &tlv.DataObject{
		TAG:  []byte{0x1F, 0x5B},
		Data: []byte{0x05},
	}
	var objects *[]tlv.DataObject = &[]tlv.DataObject{}

	*objects = append(*objects, *cardPollTimeout)
	var paymentType byte = payment.PaymentIncludeGeldKarte + payment.PrinterReady + payment.GirocardTransaction
	currency := EUR
	var tlv *tlv.Container = &tlv.Container{
		Objects: *objects,
	}
	config := AuthConfig{
		Amount:      1,
		Currency:    &currency,
		PaymentType: &paymentType,
		ExpiryDate: &ExpiryDate{
			Month: 11,
			Year:  21,
		},
		CardNumber: nil,
		TLV:        tlv,
	}
	got := compileAuthConfig(&config)
	assert.Equal(t, want, got)
}
