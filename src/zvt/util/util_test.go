package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveToFile(t *testing.T) {
	want := []byte("This are test bytes")
	filename, err := Save(&want, len(want), "Test")
	if assert.NoError(t, err) {
		got, err := Load(filename)
		if assert.NoError(t, err) {
			assert.Equal(t, want, got)
			os.Remove(filename)
		}
	}
}
