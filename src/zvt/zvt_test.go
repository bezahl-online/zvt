package zvt

import (
	"testing"

	"bezahl.online/zvt/src/zvt/payment"
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
	got := compileText([]string{"Test", "Array"})
	assert.Equal(t, want, got)
}

func TestMarshalTLV(t *testing.T) {
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
	var tlv TLV = TLV{
		Objects: *objects,
	}
	got := tlv.Marshal()
	assert.EqualValues(t, want, got)
}

func TestUnmarshalTLV(t *testing.T) {
	data := []byte{0x06, 0x0a, 0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3, 0x1f, 0x5B, 0x01, 0x05}
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
	want := TLV{
		Objects: *objects,
	}
	var got TLV = TLV{}
	err := got.Unmarshal(&data) // FIXME: test not done
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
func TestCompileLength(t *testing.T) {
	want := []byte{0x27}
	got := compileLength(39)
	if assert.Equal(t, want, got) {
		want = []byte{0x81, 0xE7}
		got = compileLength(231)
		if assert.Equal(t, want, got) {
			want = []byte{0x82, 0xD7, 0x12}
			got = compileLength(55058)
			assert.Equal(t, want, got)
		}
	}
}

func TestDecompileLength(t *testing.T) {
	want := 43
	var b []byte = []byte{0x2B, 0x1F, 0x5B}
	got, size, err := decompileLength(&b)
	if assert.NoError(t, err) {
		if assert.Equal(t, want, int(got)) && assert.Equal(t, 1, size) {
			b = []byte{0xFA, 0x1F, 0x5B}
			got, size, err = decompileLength(&b)
			if assert.Error(t, err) && assert.Equal(t, 0, size) {
				want = 250
				b = []byte{0x81, 0xFA, 0x1F, 0x5B}
				got, size, err = decompileLength(&b)
				if assert.NoError(t, err) && assert.Equal(t, 2, size) {
					if assert.Equal(t, want, int(got)) {
						want = 43981
						b = []byte{0x82, 0xab, 0xcd, 0x1F, 0x5B}
						got, size, err = decompileLength(&b)
						if assert.NoError(t, err) && assert.Equal(t, 3, size) {
							assert.Equal(t, want, int(got))

						}
					}
				}
			}
		}
	}
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
	got := compilePTConfig(&PTConfig{
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

func TestAuthData(t *testing.T) {
	want := []byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x49, 0x09, 0x97, 0x19, 0x75, 0x0E, 0x11, 0x21, 0x06, 0x04, 0x1F, 0x5B, 0x01, 0x05}
	var cardPollTimeout *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x5B},
		data: []byte{0x05},
	}
	var objects *[]DataObject = &[]DataObject{}

	*objects = append(*objects, *cardPollTimeout)
	var paymentType byte = payment.PaymentIncludeGeldKarte + payment.PrinterReady + payment.GirocardTransaction
	currency := EUR
	var tlv *TLV = &TLV{
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
	assert.EqualValues(t, want, got)
}
