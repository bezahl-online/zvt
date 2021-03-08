package bmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	data := []byte{0xF1, 0xF1, 0xF6}
	text := []byte("Das ist ein Test")
	data = append(data, text...)
	want := OBJ{
		ID:   0xF1,
		Data: text,
		Size: 2,
	}
	var got OBJ = OBJ{}
	err := got.Unmarshal(data)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
