package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompileTAG(t *testing.T) {
	want := []byte{0x49}
	data := []byte{0x49, 0x5F, 0x1A}
	got, err := DecompileTAG(&data)
	if assert.NoError(t, err) && assert.EqualValues(t, want, got) {
		want = []byte{0x1F, 0x5B}
		data = []byte{0x1F, 0x5B, 0x12, 0x45}
		got, err = DecompileTAG(&data)
		if assert.NoError(t, err) && assert.EqualValues(t, want, got) {
			data = []byte{0x1F}
			_, err = DecompileTAG(&data)
			assert.Error(t, err)
		}
	}
}
