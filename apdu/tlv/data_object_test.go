package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalDataObjects(t *testing.T) {
	want := []byte{0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3, 0x1f, 0x5B, 0x01, 0x05}
	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var cardPollTimeout *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x5B},
		Data: []byte{0x05},
	}
	var objects *[]DataObject = &[]DataObject{}

	*objects = append(*objects,
		*listOfCommands,
		*cardPollTimeout)
	got := MarshalDataObjects(objects)
	assert.EqualValues(t, want, got)
}

func TestUnmarshalDataObject(t *testing.T) {
	data := []byte{0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3}
	want := DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var got DataObject = DataObject{}
	size, err := got.Unmarshal(data)
	if assert.NoError(t, err) && assert.Equal(t, uint16(6), size) {
		if assert.Equal(t, want, got) {
			data = []byte{0x1f, 0x5B, 0x01, 0x05}
			want = DataObject{
				TAG:  []byte{0x1F, 0x5B},
				Data: []byte{0x05},
			}
			size, err := got.Unmarshal(data)
			if assert.NoError(t, err) && assert.Equal(t, uint16(4), size) {
				assert.Equal(t, want, got)
			}
		}
	}
}
