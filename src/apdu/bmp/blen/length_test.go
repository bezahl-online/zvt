package blen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthFormat(t *testing.T) {
	want := []byte{}
	l := Length{Kind: NONE, Value: uint16(123)}
	got := l.Format()
	assert.Equal(t, want, got)
	want = []byte{123}
	l = Length{Kind: BINARY, Value: uint16(123)}
	got = l.Format()
	assert.Equal(t, want, got)
	want = []byte{0xF3, 0xF8}
	l = Length{Kind: LL, Value: uint16(38)}
	got = l.Format()
	assert.Equal(t, want, got)
	want = []byte{0xF3, 0xF8, 0xF7}
	l = Length{Kind: LLL, Value: uint16(387)}
	got = l.Format()
	assert.Equal(t, want, got)
}
