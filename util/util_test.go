package util

import (
	"os"
	"testing"

	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/bezahl-online/zvt/instr"
	"github.com/stretchr/testify/assert"
)

func TestSaveToFile(t *testing.T) {
	testbytes := []byte("This are test bytes")
	want := []byte{0x06, 0xD3, 0xFF, 0x2D, 0x04}
	want = append(want, testbytes...)
	filename, err := Save(&testbytes, &instr.CtrlField{
		Class: 0x06,
		Instr: 0xD3,
		Length: blen.Length{
			Kind:  blen.BINARY,
			Size:  3,
			Value: 1069,
		},
		RawDataLength: 0,
	}, "Test")
	if assert.NoError(t, err) {
		got, err := Load(filename)
		if assert.NoError(t, err) {
			assert.Equal(t, want, got)
			os.Remove(filename)
		}
	}
}
