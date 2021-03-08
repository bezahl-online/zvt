package zvt

import (
	"testing"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/apdu/bmp/blen"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/payment"
	"bezahl.online/zvt/src/zvt/tlv"
	"bezahl.online/zvt/src/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestCommandMarshal(t *testing.T) {
	want := []byte{0x06, 0x01, 0x0a, 0x54, 0x65, 0x73,
		0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	instr := instr.CtrlField{
		Class: 0x06,
		Instr: 0x01,
		Length: blen.Length{
			Kind:  blen.BINARY,
			Value: uint16(10),
		},
	}
	c := Command{instr, apdu.DataUnit{Data: []byte("Teststring")}}
	got, err := c.Marshal()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCompileText(t *testing.T) {
	want := apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0xF1, Data: []byte{0x54, 0x65, 0x73, 0x74}},
			{ID: 0xF2, Data: []byte{0x41, 0x72, 0x72, 0x61, 0x79}},
		},
	}
	got := compileText([]string{"Test", "Array"})
	assert.Equal(t, want, got)
}

// func TestCompilePTConfig(t *testing.T) {
// 	want := []byte{0x12, 0x34, 0x56, 0x8E, 0x9, 0x78, 0x3,
// 		0x3, 0x6, 0x6, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}
// 	var listOfCommands *tlv.DataObject = &tlv.DataObject{
// 		TAG:  []byte{0x26},
// 		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
// 	}
// 	var tlvContainer *tlv.Container = &tlv.Container{
// 		Objects: []tlv.DataObject{},
// 	}
// 	tlvContainer.Objects = append(tlvContainer.Objects, *listOfCommands)
// 	config := Config{
// 		pwd:          [3]byte{0x12, 0x34, 0x56},
// 		config:       0x8E,
// 		currency:     EUR,
// 		service:      0x03,
// 		tlvContainer: tlvContainer,
// 	}
// 	got := config.CompileConfig()
// 	assert.EqualValues(t, want, got)
// }

func TestCommandUnmarshal1(t *testing.T) {
	testBytes, err := util.Load("dump/data050730027.bin")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		Instr: instr.CtrlField{
			Class: 0x04,
			Instr: 0xFF,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 2,
		},
		Data: apdu.DataUnit{
			Data:    []byte{0x0A, 0x01},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x24}, Data: testBytes[9:]},
				},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal2(t *testing.T) {
	testBytes, err := util.Load("dump/data051327012.bin")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		Instr: instr.CtrlField{
			Class: 0x06,
			Instr: 0xD3,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 0,
		},
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x1F, 0x07}, Data: testBytes[11:12]},
					{TAG: []byte{0x25}, Data: testBytes[17:]},
				},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
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
