package zvt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBytes(t *testing.T) {
	c := command{
		Class: 01,
		Inst:  02,
		Data:  []byte("Teststring"),
	}
	got := c.getBytes()
	want := []byte{0x01, 0x02, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
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

func TestCompileRegistrationData(t *testing.T) {
	// want:=[]byte{0x12,0x34,0x56,0x8E,}
}
