package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestLogOff(t *testing.T) {
	skipShort(t)
	i := instr.Map["ACK"]
	want := Command{
		CtrlField: i,
		Data:      apdu.DataUnit{},
	}
	want.CtrlField.Length.Size = 1
	err := PaymentTerminal.LogOff()
	if assert.NoError(t, err) {
		got, err := PaymentTerminal.ReadResponse()
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, *got)
		}
	}
}
