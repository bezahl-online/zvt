package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompile(t *testing.T) {
	want := [2]byte{0x49, 0}
	data := []byte{0x49, 0x5F, 0x1A}
	got, err := Decompile(&data)
	if assert.NoError(t, err) && assert.EqualValues(t, want, got) {
		want = [2]byte{0x1F, 0x5B}
		data = []byte{0x1F, 0x5B, 0x12, 0x45}
		got, err = Decompile(&data)
		if assert.NoError(t, err) && assert.EqualValues(t, want, got) {
			data = []byte{0x1F}
			_, err = Decompile(&data)
			assert.Error(t, err)
		}
	}
}
