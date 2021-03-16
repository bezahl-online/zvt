package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestDisplayText(t *testing.T) {
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
	got, err := PaymentTerminal.ReadResponse()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, *got)
	}
}
