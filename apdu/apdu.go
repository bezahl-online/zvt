package apdu

import (
	"encoding/binary"
	"fmt"

	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/bmp/blen"
	"github.com/bezahl-online/zvt/apdu/tlv"
)

// DataUnit is the structur of the APDU or APRC
type DataUnit struct {
	Data         []byte
	BMPOBJs      []bmp.OBJ
	TLVContainer tlv.Container
}

// Marshal serializes the DataUnit into bytes
func (a *DataUnit) Marshal() ([]byte, error) {
	var b []byte = []byte{}
	if a.Data != nil {
		b = append(b, a.Data...)
	}
	for _, bmp := range (*a).BMPOBJs {
		bmpBytes, err := bmp.Marshal()
		if err != nil {
			return b, err
		}
		b = append(b, bmpBytes...)
	}
	b = append(b, a.TLVContainer.Marshal()...)
	return b, nil
}

// Unmarshal unmarshals the given data
func (a *DataUnit) Unmarshal(data *[]byte) error {
	d := *data
	for {
		if len(d) < 1 {
			break
		}
		info, found := bmp.InfoMap[d[0]]
		if !found {
			break
		}
		bmpObj := bmp.OBJ{ID: d[0]}
		if len(d) < 1 {
			break
		}
		d = d[1:]
		var bmpLen uint16 = 0
		switch info.Length.Kind {
		case blen.NONE:
			bmpLen = uint16(info.Length.Value)
		case blen.BINARY: // FIXME: not tested
			if d[0] == 0xFF {
				if len(d) < 2 {
					return fmt.Errorf("wrong BMP length")
				}
				bmpLen = binary.LittleEndian.Uint16(d[1:2])
				d = d[3:]
			} else {
				d = d[1:]
			}
		case blen.LL: // FIXME: not tested
			if len(d) < 2 {
				return fmt.Errorf("wrong BMP length")
			}
			bmpLen = uint16(10*(d[1]&0x0F)) +
				uint16(d[2]&0x0F)
			d = d[2:]
		case blen.LLL: // FIXME: not tested
			if len(d) < 3 {
				return fmt.Errorf("wrong BMP length")
			}
			bmpLen = uint16(100*(d[1]&0x0F)) +
				uint16(100*(d[2]&0x0F)) +
				uint16(d[3]&0x0F)
			d = d[3:]
		}
		if len(d) < int(bmpLen) {
			dataLength := 10
			if len(d) < 10 { // FIXME: not tested
				dataLength = len(d)
			}
			return fmt.Errorf("wrong BMP length before '% X'", d[:dataLength])
		}
		bmpObj.Data = d[:bmpLen]
		d = d[bmpLen:]
		a.BMPOBJs = append(a.BMPOBJs, bmpObj)
	}
	if d[0] == tlv.BMPTLV {
		err := a.TLVContainer.Unmarshal(&d)
		if err != nil {
			return err
		}
	}
	return nil
}
