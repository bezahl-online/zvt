package aprc

// // Response is the response from the PT
// type Response struct {
// 	CCRC   byte
// 	APRC   byte
// 	Length blen.Length
// 	Data   apdu.DataUnit
// }

// // Unmarshal returns a response structure
// func (r *Response) Unmarshal(apduBytes []byte) error {
// 	var resp Response
// 	if len(apduBytes) < 3 {
// 		return fmt.Errorf("APDU less than 3 bytes long")
// 	}
// 	resp = Response{
// 		CCRC:   apduBytes[0],
// 		APRC:   apduBytes[1],
// 		Length: int(apduBytes[2]),
// 	}
// 	dataStartsAt := 3
// 	if resp.Length == 0xff {
// 		resp.Length = int(binary.LittleEndian.Uint16(apduBytes[3:4]))
// 		dataStartsAt = 5
// 	}
// 	if len(apduBytes) >= int(apduBytes[2])+dataStartsAt {
// 		resp.Data = apduBytes[dataStartsAt : apduBytes[2]+byte(dataStartsAt)]
// 	}
// 	// // instInfo := inst.InfoMaps.GetInfoMap(resp)
// 	// d := resp.Data
// 	// var dataStart, dataEnd, tagDataLength int
// 	// tagNr := d[:2]
// 	// tagInfo, found := tag.InfoMaps.GetInfoMap(tagNr)
// 	// if !found {
// 	// 	var tNr []byte = tagNr
// 	// 	if tagNr[0]&0x1F != 0x1F {
// 	// 		tNr = []byte{tagNr[0]}
// 	// 	}
// 	// 	return &resp, fmt.Errorf("TAG '% X' not found", tNr)
// 	// }
// 	// switch tagInfo.LengthType {
// 	// case tag.BINARY:
// 	// 	tagDataLength = int(d[tagInfo.TAGNrLen])
// 	// }
// 	// if tagDataLength > 0 {
// 	// 	dataStart = tagInfo.TAGNrLen + tagInfo.Length
// 	// 	dataEnd = dataStart + int(tagDataLength)
// 	// 	resp.TLV = TLV{
// 	// 		Objects: []DataObject{},
// 	// 	}
// 	// 	resp.TLV.Objects = append(resp.TLV.Objects, DataObject{
// 	// 		TAG:  []byte{0x24},
// 	// 		data: d[dataStart:dataEnd],
// 	// 	})
// 	// }
// 	return nil
// }

// // IsACK returns true if the response is in fact a ACK
// func (r *Response) IsACK() bool {
// 	return r.CCRC == 0x80 || (r.CCRC == 0x84 && r.APRC == 0)
// }

// // IsIntermediate returns true if the response is in fact a ACK
// func (r *Response) IsIntermediate() bool {
// 	return r.CCRC == 0x04 && r.APRC == 0xFF
// }

// // IsStatus returns true if the response has a status byte
// func (r *Response) IsStatus() bool {
// 	return r.CCRC == 0x04 && r.APRC == 0x0F
// }

// // GetTLV return a TLV or an error
// func (r *Response) GetTLV() (*tlv.Container, error) {
// 	var tlv tlv.Container = tlv.Container{
// 		Objects: []tlv.DataObject{},
// 	}
// 	err := tlv.Unmarshal(&r.Data)
// 	return &tlv, err
// }
