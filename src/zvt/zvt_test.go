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

func TestCompileConfigByte(t *testing.T) {
	want := ConfigByte(0x8E)
	got := ZVT.compileConfigByte([]byte{
		PaymentReceiptPrintedByECR,
		AdminReceiptPrintedByECR,
		PTSendsIntermediateStatus,
		ECRusingPrintLinesForPrintout,
	})
	assert.EqualValues(t, want, got)
}

func TestCompileServiceByte(t *testing.T) {
	want := ServiceByte(0x03)
	got := ZVT.compileServiceByte([]byte{
		ServiceMenuNOTAssignedToFunctionKey,
		DisplayTextsForCommandsAuthorisation,
	})
	assert.EqualValues(t, want, got)
}

func TestCompileTLV(t *testing.T) {
	want := []byte{0x06, 0x06, 0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3}
	got := ZVT.marshalTLV(&TLV{
		BMP:  0x06,
		data: []byte{0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3},
	})
	assert.EqualValues(t, want, got)
}

func TestCompilePTConfig(t *testing.T) {
	want := []byte{0x12, 0x34, 0x56, 0x8E, 0x9, 0x97, 0x3,
		0x3, 0x6, 0x6, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}
	got := ZVT.compilePTConfig(&PTConfig{
		pwd:      [3]byte{0x12, 0x34, 0x56},
		config:   ConfigByte(0x8E),
		currency: EUR,
		service:  ServiceByte(0x03),
		tlv: &TLV{BMP: 0x06, data: []byte{0x26, 0x04,
			0x0A, 0x02, 0x06, 0xD3}},
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
