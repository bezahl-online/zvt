package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	skipShort(t)
	err := PaymentTerminal.Status()
	assert.NoError(t, err)
}
