package tlv

// DataObject is part of a TLV
type DataObject struct {
	TAG  []byte
	Data []byte
}

// MarshalDataObjects serializes the array of TAGs with
// its optional data into a TLV data object
func MarshalDataObjects(dos *[]DataObject) []byte {
	var data []byte
	for _, obj := range *dos {
		data = append(data, obj.TAG...)
		data = append(data, CompileLength(len(obj.Data))...)
		data = append(data, obj.Data...)
	}
	return data
}

// UnmarshalDataObject retrieves a TAG with its optional data
// from a TLV data object
func UnmarshalDataObject(d []byte) (DataObject, uint16, error) {
	tag, err := DecompileTAG(&d)
	if err != nil {
		return DataObject{}, 0, err
	}
	tagLength := uint16(len(tag))
	tagLengthData := d[tagLength:]
	tagDataLength, tagLengthSize, err := DecompileLength(&tagLengthData)
	objectLength := tagLength + tagLengthSize + tagDataLength
	d = d[tagLength+tagLengthSize:]
	obj := DataObject{
		TAG:  tag,
		Data: d[:tagDataLength],
	}
	return obj, objectLength, nil
}
