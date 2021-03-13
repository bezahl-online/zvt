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
	want := []byte{0x80, 0}
	want = append(want, testbytes...)
	filename, err := Save(&testbytes, &instr.CtrlField{
		Class: 0x80,
		Instr: 00,
		Length: blen.Length{
			Kind:  0,
			Size:  0,
			Value: 0,
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
