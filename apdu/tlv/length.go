package tlv

import (
	"encoding/binary"
	"fmt"
)

// CompileLength constructs the length byte(s) according
// to the ZVT protocol lengths for TLV data objects
func CompileLength(len int) []byte {
	var length []byte = []byte{0}
	if len > 255 {
		length[0] = 0x82
		var l []byte = []byte{0, 0}
		binary.BigEndian.PutUint16(l, uint16(len))
		length = append(length, l...)
	} else if len > 127 {
		length[0] = 0x81
		length = append(length, byte(len))
	} else {
		length[0] = byte(len)
	}
	return length
}

// DecompileLength retrieves the length from the big edian
// coded length in the TLV data object data structure byte(s)
// according to the ZVT protocol
// returns (dataLength(uint16), tagLengthSize(uint16) error)
func DecompileLength(data *[]byte) (uint16, uint16, error) {
	l := *data
	if l[0]&0x80 == 0x80 {
		if l[0] == 0x82 && len(l) >= 3 {
			return binary.BigEndian.Uint16(l[1:3]), 3, nil
		} else if l[0] == 0x81 && len(l) >= 2 {
			return uint16(l[1]), 2, nil
		}
		return 0, 0, fmt.Errorf("invalid value")
	}
	return uint16(l[0] & 0x7F), 1, nil
}
