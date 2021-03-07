package instr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	want := []byte{0x06, 0x00, 0x0D}
	ctrlField := Map["Registration"]
	got := ctrlField.Marshal(uint16(13))
	assert.Equal(t, want, got)
}
