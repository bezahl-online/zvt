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

func TestUnmarshal(t *testing.T) {
	want := uint16(123)
	l := Length{
		Kind:  BINARY,
		Value: 0,
	}
	err := l.Unmarshal([]byte{123})
	if assert.NoError(t, err) {
		got := l.Value
		assert.Equal(t, want, got)
	}
	want = uint16(23545)
	l = Length{
		Kind:  BINARY,
		Value: 0,
	}
	err = l.Unmarshal([]byte{0xff, 0xf9, 0x5b})
	if assert.NoError(t, err) {
		got := l.Value
		assert.Equal(t, want, got)
	}
	want = uint16(125)
	l = Length{
		Kind:  LLL,
		Value: 0,
	}
	err = l.Unmarshal([]byte{0xf1, 0xf2, 0xf5})
	if assert.NoError(t, err) {
		got := l.Value
		assert.Equal(t, want, got)
	}
	want = uint16(38)
	l = Length{
		Kind:  LL,
		Value: 0,
	}
	err = l.Unmarshal([]byte{0xf3, 0xf8})
	if assert.NoError(t, err) {
		got := l.Value
		assert.Equal(t, want, got)
	}
}
