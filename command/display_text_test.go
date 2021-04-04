package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestDisplayText(t *testing.T) {
	skipShort(t)
	i := instr.Map["ACK"]
	want := Command{
		CtrlField: i,
		Data:      apdu.DataUnit{},
	}
	want.CtrlField.Length.Size = 1
	err := PaymentTerminal.DisplayText([]string{
		"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	if assert.NoError(t, err) {
		got, err := PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, *got)
		}
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
