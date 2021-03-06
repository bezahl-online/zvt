package length

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthFormat(t *testing.T) {
	want := []byte{}
	got := Format(123, NONE)
	assert.Equal(t, want, got)
	want = []byte{123}
	got = Format(123, BINARY)
	assert.Equal(t, want, got)
	want = []byte{0xF3, 0xF8}
	got = Format(38, LL)
	assert.Equal(t, want, got)
	want = []byte{0xF3, 0xF8, 0xF7}
	got = Format(387, LLL)
	assert.Equal(t, want, got)
}
