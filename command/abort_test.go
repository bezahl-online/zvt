package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
	skipShort(t)
	err := PaymentTerminal.Abort()
	assert.NoError(t, err)
}
