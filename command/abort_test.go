package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
	err := PaymentTerminal.Abort()
	assert.NoError(t, err)
}
