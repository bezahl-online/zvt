package instr

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	want := []byte{0x06, 0x00, 0x0D}
	ctrlField := Map["Registration"]
	got, err := ctrlField.Marshal(uint16(13))
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
}

func TestFind(t *testing.T) {
	want := CtrlField{
		Class: 0x06,
		Instr: 0x00,
		Length: blen.Length{
			Kind:  blen.BINARY,
			Value: 0,
		},
	}
	d := []byte{0x06, 0x00}
	got := Find(&d)
	assert.Equal(t, want, *got)
}
