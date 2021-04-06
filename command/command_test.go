package command

import (
	"testing"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"github.com/stretchr/testify/assert"
)

func TestCommandMarshal(t *testing.T) {
	// want := []byte{0x06, 0x01, 0x0a, 0x54, 0x65, 0x73,
	// 	0x74, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	instr := instr.CtrlField{
		Class: 0x06,
		Instr: 0x01,
		Length: blen.Length{
			Kind:  blen.BINARY,
			Value: uint16(10),
		},
	}
	want := []byte{0x06, 0x01, 0x17, 0x04, 0, 0, 0, 1, 0, 0, 0x05, 0x02, 0x06, 0x0C, 0x06, 0x06, 0x26, 0x4, 0xa, 0x2, 0x6, 0xd3, 0x1F, 0x04, 0x01, 0x02}
	var apdu apdu.DataUnit = apdu.DataUnit{
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
	c := Command{
		CtrlField: instr,
		Data:      apdu,
	}
	got, err := c.Marshal()
	if assert.NoError(t, err) {
		if assert.EqualValues(t, want, got) {
			c.Data.BMPOBJs[0].ID = 0xFF
			_, err := c.Marshal()
			assert.Error(t, err)
		}
	}
}

func TestCommandMarshal0(t *testing.T) {
	instr := instr.Map["ACK"]
	want := []byte{0x80, 0x00, 0x00}
	c := Command{
		CtrlField: instr,
	}
	got, err := c.Marshal()
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal1(t *testing.T) {
	testBytes, err := util.Load("testdata/data050730027.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0xFF,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 2,
		},
		Data: apdu.DataUnit{
			Data:    []byte{0x0A, 0x01},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x24}, Data: testBytes[9:]},
				},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal2(t *testing.T) {
	testBytes, err := util.Load("testdata/data051327012.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x06,
			Instr: 0xD3,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 0,
		},
		Data: apdu.DataUnit{
			Data:    []byte{},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x1F, 0x07}, Data: testBytes[11:12]},
					{TAG: []byte{0x25}, Data: testBytes[17:]},
				},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal3(t *testing.T) {
	testBytes, err := util.Load("testdata/data081537024.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0xFF,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 2,
		},
		Data: apdu.DataUnit{
			Data:    []byte{0x0A, 0x01},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{
					{TAG: []byte{0x24}, Data: testBytes[9:]},
				},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal4(t *testing.T) {
	testBytes, err := util.Load("testdata/data081537039.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0x0F,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 0,
		},
		Data: apdu.DataUnit{
			Data: []byte{},
			BMPOBJs: []bmp.OBJ{
				{ID: 0x27, Data: []byte{0x6C}, Size: 2},
				{ID: 0x29, Data: []byte{0x29, 0, 0x10, 0x6}, Size: 5},
			},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
func TestCommandUnmarshal5(t *testing.T) {
	testBytes, err := util.Load("testdata/data081621025.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0xFF,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 2,
		},
		Data: apdu.DataUnit{
			Data:    []byte{0x0E},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

var Objects []bmp.OBJ = []bmp.OBJ{
	{ID: 0x27, Size: 2, Data: []uint8{0x0}},
	{ID: 0x4, Size: 7, Data: []uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x1}},

	{ID: 0xb, Size: 4, Data: []uint8{0x0, 0x0, 0x22}},
	{ID: 0xc, Size: 4, Data: []uint8{0x16, 0x49, 0x23}},
	{ID: 0xd, Size: 3, Data: []uint8{0x3, 0x8}},
	{ID: 0x17, Size: 3, Data: []uint8{0x0, 0x1}},
	{ID: 0x19, Size: 2, Data: []uint8{0x70}},
	{ID: 0x22, Size: 11, Data: []uint8{0xee, 0xee, 0xee, 0xee, 0xee,
		0xee, 0x57, 0x26}},
	{ID: 0x29, Size: 5, Data: []uint8{0x29, 0x0, 0x10, 0x6}},
	{ID: 0x2a, Size: 16, Data: []uint8{0x31, 0x30, 0x30, 0x37, 0x36,
		0x34, 0x39, 0x39, 0x32, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}},
	{ID: 0x3b, Size: 9, Data: []uint8{0x32, 0x39, 0x31, 0x36, 0x37,
		0x35, 0x0, 0x0}},
	{ID: 0x3c, Size: 18, Data: []uint8{0x47, 0x45, 0x4e, 0x2e, 0x4e,
		0x52, 0x2e, 0x3a, 0x32, 0x39, 0x31, 0x36, 0x37, 0x35}},
	{ID: 0x49, Size: 3, Data: []uint8{0x9, 0x78}},
	{ID: 0x87, Size: 3, Data: []uint8{0x0, 0x22}},
	{ID: 0x88, Size: 4, Data: []uint8{0x0, 0x0, 0x22}},
	{ID: 0x8a, Size: 2, Data: []uint8{0x2e}},
	{ID: 0x8b, Size: 20, Data: []uint8{0x44, 0x65, 0x62, 0x69, 0x74,
		0x20, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x63, 0x61, 0x72, 0x64, 0x0}}}

func TestCommandUnmarshal6(t *testing.T) {
	testBytes, err := util.Load("testdata/data081649050.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{CtrlField: instr.CtrlField{
		Class: 0x4, Instr: 0xf, Length: blen.Length{Kind: 0x1, Size: 0x0, Value: 0x0}, RawDataLength: 0},
		Data: apdu.DataUnit{Data: []uint8{}, BMPOBJs: Objects, TLVContainer: tlv.Container{Objects: []tlv.DataObject{
			{TAG: []uint8{0x24}, Data: []uint8{0x7, 0xe, 0x47, 0x45,
				0x4e, 0x2e, 0x4e, 0x52, 0x2e, 0x3a, 0x32, 0x39, 0x31, 0x36, 0x37, 0x35}},
			{TAG: []uint8{0x43}, Data: []uint8{0xa0, 0x0, 0x0, 0x0,
				0x4, 0x10, 0x10}},
			{TAG: []uint8{0x46}, Data: []uint8{0x41, 0x30, 0x30, 0x30,
				0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x31, 0x30, 0x31, 0x30, 0x20, 0x35,
				0x46, 0x32, 0x36, 0x30, 0x45, 0x46, 0x43, 0x35, 0x44, 0x32, 0x38, 0x45,
				0x34, 0x35, 0x31}},
			{TAG: []uint8{0x47}, Data: []uint8{0x41, 0x30, 0x30, 0x30,
				0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x31, 0x30, 0x31, 0x30, 0x20, 0x35,
				0x46, 0x32, 0x36, 0x30, 0x45, 0x46, 0x43, 0x35, 0x44, 0x32, 0x38, 0x45,
				0x34, 0x35, 0x31}},
			{TAG: []uint8{0x1f, 0x10}, Data: []uint8{0x0}},
			{TAG: []uint8{0x1f, 0x12}, Data: []uint8{0x2}}}}}}

	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestFromOBJs(t *testing.T) {
	want := AuthResultData{
		Amount:     1,
		ReceiptNr:  22,
		TurnoverNr: 22,
		TraceNr:    22,
		Date:       "0308",
		Time:       "164923",
		TID:        "29001006",
		VU:         "100764992",
		AID:        "291675",
		Card: CardData{
			Name:  "Debit Mastercard",
			Type:  46,
			PAN:   "XXXX XXXX XXXX 5726",
			Tech:  0,
			SeqNr: 1,
		},
	}
	var got AuthResultData = AuthResultData{}
	got.FromOBJs(Objects)
	assert.EqualValues(t, want, got)
}

func TestCommandUnmarshal7(t *testing.T) {
	testBytes, err := util.Load("testdata/data091355012.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x06,
			Instr: 0x1E,
			Length: blen.Length{
				Kind: blen.BINARY,
			},
			RawDataLength: 1,
		},
		Data: apdu.DataUnit{
			Data:    []byte{0x6C},
			BMPOBJs: []bmp.OBJ{},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}

func TestCommandUnmarshal8(t *testing.T) {
	testBytes, err := util.Load("testdata/1617181236803PT.hex")
	if !assert.NoError(t, err) {
		return
	}
	want := Command{
		CtrlField: instr.CtrlField{
			Class: 0x04,
			Instr: 0x0F,
			Length: blen.Length{
				Kind:  blen.BINARY,
				Size:  0,
				Value: 0,
			},
			RawDataLength: 0,
		},
		Data: apdu.DataUnit{
			Data: []byte{},
			BMPOBJs: []bmp.OBJ{
				{ID: 0x27, Size: 2, Data: []byte{0}},
				{ID: 0x04, Size: 7, Data: []byte{0, 0, 0, 0, 0x26, 0}},
				{ID: 0x0C, Size: 4, Data: []byte{0x10, 0x59, 0x37}},
				{ID: 0x0D, Size: 3, Data: []byte{0x03, 0x31}},
				{ID: 0x60, Size: 57, Data: []byte{0, 0x68, 0, 0x71, 0x04,
					0, 0, 0, 0, 0x26, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			},
			TLVContainer: tlv.Container{
				Objects: []tlv.DataObject{},
			},
		},
	}
	var got Command = Command{}
	err = got.Unmarshal(&testBytes)
	if assert.NoError(t, err) {
		assert.EqualValues(t, want, got)
	}
}
