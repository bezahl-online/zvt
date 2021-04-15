package command

import (
	"testing"

	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	skipShort(t) // start
	i := instr.Map["ACK"]
	i.Length.Size = 1
	err := PaymentTerminal.Register()
	assert.NoError(t, err)
}
