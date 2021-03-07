package bmp

import (
	"fmt"
)

// OBJ is
type OBJ struct {
	ID   byte
	Data []byte
}

// Marshal serializes the BMP with its specific length field
func (o *OBJ) Marshal() ([]byte, error) {
	var b []byte = []byte{o.ID}
	info, found := InfoMap[o.ID]
	if !found {
		return b, fmt.Errorf("BMP with ID % X not found", o.ID)
	}
	info.Length.Value = uint16(len(o.Data))
	b = append(b, info.Length.Format()...)
	b = append(b, o.Data...)
	return b, nil
}
