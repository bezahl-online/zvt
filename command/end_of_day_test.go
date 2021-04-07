package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndOfDay(t *testing.T) {
	skipShort(t)
	// t.Cleanup(func() {
	// 	PaymentTerminal.Abort()
	// })
	err := PaymentTerminal.EndOfDay()
	assert.NoError(t, err)

}
