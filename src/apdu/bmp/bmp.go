package bmp

import (
	"fmt"
)

// OBJ is
type OBJ struct {
	ID   byte
	Size int // length of the Lenght
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

// Unmarshal it
func (o *OBJ) Unmarshal(data []byte) error {
	var err error
	if len(data) < 1 {
		return fmt.Errorf("no data")
	}
	o.ID = data[0]
	info, found := InfoMap[o.ID]
	if !found {
		return fmt.Errorf("BMP % X not found", data[0])
	}
	err = info.Length.Unmarshal(data[1:])
	if err == nil {
		start := info.Length.Size + 1
		o.Size = int(info.Length.Size)
		o.Data = data[start : info.Length.Value+uint16(start)]
	}
	return err
}
