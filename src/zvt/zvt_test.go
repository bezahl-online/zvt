package zvt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBytes(t *testing.T) {
	c := Command{
		Class: 01,
		Inst:  02,
		Data:  []byte("Teststring"),
	}
	got := c.getBytes()
	want := []byte{0x01, 0x02, 0x0a, 0x54, 0x65, 0x73,
		0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	assert.Equal(t, want, got)
}

func TestCompileText(t *testing.T) {
	want := []byte{0xf1, 0xf0, 0xf4, 0x54, 0x65, 0x73,
		0x74, 0xf2, 0xf0, 0xf5, 0x41, 0x72, 0x72, 0x61, 0x79}
	got := ZVT.compileText([]string{"Test", "Array"})
	assert.Equal(t, want, got)
}

func TestCompileTLV(t *testing.T) {
	want := []byte{0x06, 0x0a, 0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3, 0x1f, 0x5B, 0x01, 0x05}
	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var cardPollTimeout *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x5B},
		data: []byte{0x05},
	}
	var objects *[]DataObject = &[]DataObject{}

	*objects = append(*objects,
		*listOfCommands,
		*cardPollTimeout)
	got := ZVT.marshalTLV(objects)
	assert.EqualValues(t, want, got)
}

func TestCompilePTConfig(t *testing.T) {
	want := []byte{0x12, 0x34, 0x56, 0x8E, 0x9, 0x97, 0x3,
		0x3, 0x6, 0x6, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}
	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlv *TLV = &TLV{
		Objects: []DataObject{},
	}
	tlv.Objects = append(tlv.Objects, *listOfCommands)
	got := ZVT.compilePTConfig(&PTConfig{
		pwd:      [3]byte{0x12, 0x34, 0x56},
		config:   0x8E,
		currency: EUR,
		service:  0x03,
		tlv:      tlv,
	})
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
		assert.EqualValues(t, want, got)
	}
}

// func TestAuthData(t *testing.T){
// 	want := []byte{0x040, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x19, 0x74, 0x0E,0x11,0x21,0x06,0x03,0x1F,0x5B,0x15}
// 	data:=AuthData{
// 		Amount:      0,
// 		Currency:    new(int),
// 		PaymentType: new(byte),
// 		ExpiryDate:  &ExpiryDate{},
// 		CardNumber:  new(int),
// 		TLV:         &TLV{
// 			BMP:  ,
// 			data: []byte{},
// 		},
// 	}
// 	got, err := ZVT.marshalAuthData(data)
// 	if assert.NoError(t, err) {
// 		assert.EqualValues(t, want, got)
// 	}
//
