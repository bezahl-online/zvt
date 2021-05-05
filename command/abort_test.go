package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
	if skipShort(t) {
		return
	}
	err := PaymentTerminal.Abort()
	assert.NoError(t, err)
}
