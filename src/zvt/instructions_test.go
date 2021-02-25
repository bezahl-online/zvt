package zvt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayText(t *testing.T) {
	err := ZVT.DisplayText([]string{
		"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	assert.NoError(t, err)
}
