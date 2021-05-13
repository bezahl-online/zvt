package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/stretchr/testify/assert"
)

func TestDisplayText(t *testing.T) {
	if skipShort(t) {
		return
	}
	PaymentTerminal.DisplayText([]string{
		// "Da steh ich nun,",
		// "ich armer Tor,",
		// "Und bin so klug",
		// "als wie zuvor.",
		"Greisslomat e.U.",
		"Dorf 122",
		"6645 Vorderhornbach",
	})

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
