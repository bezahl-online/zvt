package zvt

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/zvt/length"
	"bezahl.online/zvt/src/zvt/util"
	"github.com/albenik/bcd"
)

// ZVT represents the driver
var ZVT PT
var zvtACK Response = Response{
	CCRC:   0x80,
	APRC:   0,
	Length: 0,
	Data:   []byte{},
}

func init() {
	var pt PT = PT{
		lock: &sync.RWMutex{},
		conn: nil,
	}
	err := pt.Open()
	if err != nil {
		fmt.Println(err.Error())
	}
	ZVT = pt
}

// Open opens a connection to the PT
func (p *PT) Open() error {
	var url string
	if len(os.Getenv("ZVT_URL")) > 0 {
		url = os.Getenv("ZVT_URL")
	} else {
		return fmt.Errorf("Please set environment variabel ZVT_URL")
	}
	var err error
	p.conn, err = net.Dial("tcp", url)
	if err != nil {
		return err
	}
	p.conn.SetDeadline(time.Now().Add(5 * time.Second))
	return nil
}

// convert command structure to byte array
func (c *Command) compile() ([]byte, error) {
	var b []byte = []byte{
		c.Class,
		c.Inst,
	}
	data, err := c.Data.Marshal()
	if err != nil {
		return b, err
	}
	b = append(b, length.Format(uint16(len(data)), length.BINARY)...)
	b = append(b, data...)
	return b, nil
}

// SendACK send ACK and return the response or error
func (p *PT) SendACK(timeout time.Duration) (*Response, error) {
	ack := zvtACK.Marshal()
	nr, err := p.conn.Write(ack)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ECR => PT (%3d):% X\n", nr, ack)
	return p.readResponse(timeout)
}

func (p *PT) send(c Command) (*Response, error) {
	var err error
	b, err := c.compile()
	if err != nil {
		return nil, err
	}
	fmt.Printf("% X", b)
	_, err = p.conn.Write(b)
	if err != nil {
		return nil, err
	}

	var resp *Response
	resp, err = p.readResponse(5 * time.Second)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (p *PT) readResponse(timeout time.Duration) (*Response, error) {
	var resp *Response
	var err error
	var readBuf []byte = make([]byte, 1024)
	p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err := p.conn.Read(readBuf)
	util.Save(&readBuf, nr)
	if err != nil {
		return resp, err
	}
	fmt.Printf("PT => ECR (%3d):% X\n", nr, readBuf[:nr])
	resp, err = p.unmarshalAPDU(readBuf[:nr])
	return resp, err
}

func compileText(textarray []string) apdu.DataUnit {
	dataUnit := apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{},
	}
	for i, text := range textarray {
		bmp := bmp.OBJ{
			ID:   0xf0 + byte(i+1),
			Data: []byte{},
		}
		bmp.Data = append(bmp.Data, []byte(text)...) // append bytes of text line
		dataUnit.BMPOBJs = append(dataUnit.BMPOBJs, bmp)
	}
	return dataUnit
}

func compileLL(l uint8) []byte {
	var b []byte = make([]byte, 2)
	lz := uint8(l / 10)           // value of tens
	le := uint8(l - uint8(10*lz)) // value of unit position
	b[0] = 0xF0 + lz              // code into 0xFx (tens) (BCD)
	b[1] = 0xF0 + le              // code into 0xFy (unit) (BCD)
	return b
}

func (p *PT) unmarshalAPDU(apduBytes []byte) (*Response, error) {
	var resp Response
	if len(apduBytes) < 3 {
		return &resp, fmt.Errorf("APDU less than 3 bytes long")
	}
	resp = Response{
		CCRC:   apduBytes[0],
		APRC:   apduBytes[1],
		Length: int(apduBytes[2]),
	}
	dataStartsAt := 3
	if resp.Length == 0xff {
		resp.Length = int(binary.LittleEndian.Uint16(apduBytes[3:4]))
		dataStartsAt = 5
	}
	if len(apduBytes) >= int(apduBytes[2])+dataStartsAt {
		resp.Data = apduBytes[dataStartsAt : apduBytes[2]+byte(dataStartsAt)]
	}
	// // instInfo := inst.InfoMaps.GetInfoMap(resp)
	// d := resp.Data
	// var dataStart, dataEnd, tagDataLength int
	// tagNr := d[:2]
	// tagInfo, found := tag.InfoMaps.GetInfoMap(tagNr)
	// if !found {
	// 	var tNr []byte = tagNr
	// 	if tagNr[0]&0x1F != 0x1F {
	// 		tNr = []byte{tagNr[0]}
	// 	}
	// 	return &resp, fmt.Errorf("TAG '% X' not found", tNr)
	// }
	// switch tagInfo.LengthType {
	// case tag.BINARY:
	// 	tagDataLength = int(d[tagInfo.TAGNrLen])
	// }
	// if tagDataLength > 0 {
	// 	dataStart = tagInfo.TAGNrLen + tagInfo.Length
	// 	dataEnd = dataStart + int(tagDataLength)
	// 	resp.TLV = TLV{
	// 		Objects: []DataObject{},
	// 	}
	// 	resp.TLV.Objects = append(resp.TLV.Objects, DataObject{
	// 		TAG:  []byte{0x24},
	// 		data: d[dataStart:dataEnd],
	// 	})
	// }
	return &resp, nil
}

func compileAuthConfig(c *AuthConfig) apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x49, Data: bcd.FromUint16(uint16(*c.Currency))},
			{ID: 0x19, Data: []byte{*c.PaymentType}},
		},
	}
}
