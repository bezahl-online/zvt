package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorisation(t *testing.T) {
	t.Cleanup(func() {
		// PaymentTerminal.Abort()
	})
	config := &AuthConfig{
		Amount: 1,
	}
	err := PaymentTerminal.Authorisation(config)
	if assert.NoError(t, err) {
	}
}

func TestFormatPAN(t *testing.T) {
	want := "XXXX XXXX XXXX 5726"
	data := []uint8{0xee, 0xee, 0xee, 0xee, 0xee,
		0xee, 0x57, 0x26}
	got := formatPAN(data)
	assert.Equal(t, want, got)
	want = "XXXX XXXX XXXX 572"
	data = []uint8{0xee, 0xee, 0xee, 0xee, 0xee,
		0xee, 0x57, 0x2F}
	got = formatPAN(data)
	assert.Equal(t, want, got)
}
