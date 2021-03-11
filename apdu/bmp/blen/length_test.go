package blen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthFormat(t *testing.T) {
	var l Length
	var got []byte
	var err error
	l = Length{Kind: 255, Value: uint16(123)}
	got, err = l.Format()
	assert.Error(t, err)
	want := []byte{}
	l = Length{Kind: NONE, Value: uint16(123)}
	got, err = l.Format()
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
	want = []byte{123}
	l = Length{Kind: BINARY, Value: uint16(123)}
	got, err = l.Format()
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
	want = []byte{0xff, 0x6F, 0x02}
	l = Length{Kind: BINARY, Value: uint16(623)}
	got, err = l.Format()
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
	want = []byte{0xF3, 0xF8}
	l = Length{Kind: LL, Value: uint16(38)}
	got, err = l.Format()
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
	want = []byte{0xF3, 0xF8, 0xF7}
	l = Length{Kind: LLL, Value: uint16(387)}
	got, err = l.Format()
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
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
