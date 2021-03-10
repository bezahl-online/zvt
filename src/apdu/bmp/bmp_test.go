package bmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	want := []byte{0xF1, 0xF1, 0xF6}
	text := []byte("Das ist ein Test")
	want = append(want, text...)
	obj := OBJ{
		ID:   0xF1,
		Data: text,
		Size: 19,
	}
	got, err := obj.Marshal()
	if assert.NoError(t, err) {
		if assert.EqualValues(t, want, got) {
			obj := OBJ{
				ID: 0xf9,
			}
			_, err := obj.Marshal()
			assert.Error(t, err)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte{0xF1, 0xF1, 0xF6}
	text := []byte("Das ist ein Test")
	data = append(data, text...)
	want := OBJ{
		ID:   0xF1,
		Data: text,
		Size: 19,
	}
	var got OBJ = OBJ{}
	err := got.Unmarshal(data)
	if assert.NoError(t, err) {
		if assert.EqualValues(t, want, got) {
			err := got.Unmarshal([]byte{0xf9, 0xf1, 0xF6})
			if assert.Error(t, err) {
				err = got.Unmarshal([]byte{})
				assert.Error(t, err)
			}
		}
	}
}
