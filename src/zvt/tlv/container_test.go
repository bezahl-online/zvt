package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalTLV(t *testing.T) {
	want := []byte{0x06, 0x0a, 0x26, 0x04, 0x0A, 0x02, 0x06, 0xD3, 0x1f, 0x5B, 0x01, 0x05}
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
	var tlvContainer Container = Container{
		Objects: *objects,
	}
	got := tlvContainer.Marshal()
	assert.EqualValues(t, want, got)
}

func TestUnmarshalTLV(t *testing.T) {

	data := []byte{0x06, 0x10, 0x26, 0x04, 0x0A, 0x02, 0x06,
		0xD3, 0x1f, 0x5B, 0x01, 0x05, 0x45, 0x04, 0x02, 0x02, 0, 0}
	var listOfCommands *DataObject = &DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var cardPollTimeout *DataObject = &DataObject{
		TAG:  []byte{0x1F, 0x5B},
		Data: []byte{0x05},
	}
	var receiptParameter *DataObject = &DataObject{
		TAG:  []byte{0x45},
		Data: []byte{0x02, 0x02, 0, 0},
	}
	var objects *[]DataObject = &[]DataObject{}
	*objects = append(*objects,
		*listOfCommands,
		*cardPollTimeout,
		*receiptParameter,
	)
	want := Container{
		Objects: *objects,
	}
	var got Container = Container{}
	err := got.Unmarshal(&data)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
