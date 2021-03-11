package tlv

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/bezahl-online/zvt/apdu/tlv/tag"
)

// DataObject is part of a TLV
type DataObject struct {
	TAG  []byte
	Data []byte
}

// MarshalDataObjects serializes the array objects
func MarshalDataObjects(dos *[]DataObject) []byte {
	var data []byte
	for _, obj := range *dos {
		data = append(data, obj.Marshal()...)
	}
	return data
}

// Marshal serializes a TAG with its optional data
// into a bytes of a TLV data object
func (obj *DataObject) Marshal() []byte {
	var data []byte
	data = append(data, obj.TAG...)
	data = append(data, CompileLength(len(obj.Data))...)
	data = append(data, obj.Data...)
	return data
}

// Unmarshal retrieves a TAG with its optional data
// from the bytes of a TLV data object
func (obj *DataObject) Unmarshal(d []byte) (uint16, error) {
	tagNr, err := tag.Decompile(&d)
	if err != nil {
		return 0, err
	}
	info, found := tag.InfoMap[tagNr]
	if !found {
		return 0, fmt.Errorf("TAG '%04X' not found", tagNr)
	}
	tagLengthSize := uint16(0)
	tagDataLength := uint16(info.Length)
	objectLength := uint16(info.TAGNrLen) + tagLengthSize + tagDataLength
	tagLengthData := d[info.TAGNrLen:]
	switch info.LengthType {
	case blen.BINARY:
		tagDataLength, tagLengthSize, err = DecompileLength(&tagLengthData)
		if err != nil {
			return 0, err
		}
		objectLength = uint16(info.TAGNrLen) + tagLengthSize + tagDataLength
		// case blen.BCD:
		// 	tagDataLength = uint16(info.Length)
	}
	d = d[info.TAGNrLen+int(tagLengthSize):]
	(*obj).Data = d[:tagDataLength]
	(*obj).TAG = tagNr[:info.TAGNrLen]
	return objectLength, nil
}
