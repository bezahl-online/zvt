package bmp

import (
	"fmt"

	"bezahl.online/zvt/src/zvt/length"
)

// OBJ is
type OBJ struct {
	ID   byte
	Data []byte
}

// Marshal serializes the BMP with its specific length field
func (o *OBJ) Marshal() ([]byte, error) {
	var b []byte = []byte{o.ID}
	dataLen := uint16(len(o.Data))
	info, found := InfoMap[o.ID]
	if !found {
		return b, fmt.Errorf("BMP with ID % X not found", o.ID)
	}
	b = append(b, length.Format(dataLen, info.LengthType)...)
	b = append(b, o.Data...)
	return b, nil
}
