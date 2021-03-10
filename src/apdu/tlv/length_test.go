package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileLength(t *testing.T) {
	want := []byte{0x27}
	got := CompileLength(39)
	if assert.Equal(t, want, got) {
		want = []byte{0x81, 0xE7}
		got = CompileLength(231)
		if assert.Equal(t, want, got) {
			want = []byte{0x82, 0xD7, 0x12}
			got = CompileLength(55058)
			assert.Equal(t, want, got)
		}
	}
}

func TestDecompileLength(t *testing.T) {
	want := uint16(43)
	var b []byte = []byte{0x2B, 0x1F, 0x5B}
	got, size, err := DecompileLength(&b)
	if assert.NoError(t, err) {
		if assert.Equal(t, want, got) && assert.Equal(t, uint16(1), size) {
			b = []byte{0xFA, 0x1F, 0x5B}
			got, size, err = DecompileLength(&b)
			if assert.Error(t, err) && assert.Equal(t, uint16(0), size) {
				want = uint16(250)
				b = []byte{0x81, 0xFA, 0x1F, 0x5B}
				got, size, err = DecompileLength(&b)
				if assert.NoError(t, err) && assert.Equal(t, uint16(2), size) {
					if assert.Equal(t, want, got) {
						want = 43981
						b = []byte{0x82, 0xab, 0xcd, 0x1F, 0x5B}
						got, size, err = DecompileLength(&b)
						if assert.NoError(t, err) && assert.Equal(t, uint16(3), size) {
							assert.Equal(t, want, got)

						}
					}
				}
			}
		}
	}
}
