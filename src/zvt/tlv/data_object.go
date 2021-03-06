package tlv

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
	tag, err := DecompileTAG(&d)
	if err != nil {
		return 0, err
	}
	tagLength := uint16(len(tag))
	tagLengthData := d[tagLength:]
	tagDataLength, tagLengthSize, err := DecompileLength(&tagLengthData)
	objectLength := tagLength + tagLengthSize + tagDataLength
	d = d[tagLength+tagLengthSize:]
	(*obj).TAG = tag
	(*obj).Data = d[:tagDataLength]
	return objectLength, nil
}
