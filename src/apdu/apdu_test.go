package apdu

import (
	"testing"

	"github.com/bezahl-online/zvt/src/apdu/bmp"
	"github.com/bezahl-online/zvt/src/apdu/tlv"
	"github.com/bezahl-online/zvt/src/zvt/config"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	want := []byte{0x04, 0, 0, 0, 1, 0, 0, 0x05, 0x02, 0x06, 0x0C, 0x06, 0x06, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3, 0x1F, 0x04, 0x01, 0x02}
	var apdu DataUnit = DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: []byte{0, 0, 0, 1, 0, 0}},
			{ID: 0x05, Data: []byte{2}},
		},
		TLVContainer: tlv.Container{
			Objects: []tlv.DataObject{
				{TAG: []byte{0x06}, Data: []byte{0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}},
				{TAG: []byte{0x1F, 0x04}, Data: []byte{0x02}},
			},
		},
	}
	got, err := apdu.Marshal()
	if assert.NoError(t, err) && assert.EqualValues(t, want, got) {
		want = []byte{0x8e, 0x3, 0x6, 0xc, 0x1f, 0x73, 0x3, 0x0, 0x0, 0x0, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3}
		configByte := byte(config.PaymentReceiptPrintedByECR +
			config.AdminReceiptPrintedByECR +
			config.PTSendsIntermediateStatus +
			config.ECRusingPrintLinesForPrintout)
		serviceByte := byte(config.ServiceMenuNOTAssignedToFunctionKey +
			config.DisplayTextsForCommandsAuthorisation)
		var apdu DataUnit = DataUnit{
			Data: []byte{configByte, serviceByte},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x1F, 0x73}, Data: []byte{0, 0, 0}},
					{TAG: []byte{0x26}, Data: []byte{0x0A, 0x02, 0x06, 0xD3}},
				},
			},
		}
		got, err = apdu.Marshal()
		if assert.NoError(t, err) {
			assert.EqualValues(t, want, got)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte{0x19, 0x01, 0x29, 0x29, 0x00, 0x10, 0x06, 0x49, 0x09, 0x78, 0x04, 0, 0, 0, 1, 0, 0, 0x06, 0x04, 0x1F, 0x04, 0x01, 0x02}
	want := DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x19, Data: []byte{0x01}},
			{ID: 0x29, Data: []byte{0x29, 0x00, 0x10, 0x06}},
			{ID: 0x49, Data: []byte{0x09, 0x78}},
			{ID: 0x04, Data: []byte{0, 0, 0, 1, 0, 0}},
		},
		TLVContainer: tlv.Container{
			Objects: []tlv.DataObject{
				{TAG: []byte{0x1F, 0x04}, Data: []byte{0x02}},
			},
		},
	}
	var got DataUnit = DataUnit{}
	err := got.Unmarshal(&data)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
