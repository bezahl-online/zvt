package zvt

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

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
func (c Command) compile() []byte {
	var b []byte = []byte{
		c.Class,
		c.Inst,
		c.Length,
	}
	b[2] = byte(len(c.Data))
	b = append(b, c.Data...)
	return b
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
	nr, err := p.conn.Write(c.compile())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ECR => PT (%3d):% X\n", nr, c.compile())
	// fmt.Println(strings.ReplaceAll(fmt.Sprintf("%q\n", fmt.Sprintf("% #x", c.getBytes())), " ", ","))
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

func compileText(textarray []string) []byte {
	var t []byte = []byte{}
	for i, text := range textarray {
		fmt.Println(text)
		te := []byte{0xf0 + byte(i+1), 0x0, 0x0}
		t1 := []byte(text)
		l := uint8(len(t1))
		lz := uint8(l / 10)
		le := uint8(l - uint8(10*lz))
		te[1] = 0xf0 + lz
		te[2] = 0xf0 + le
		te = append(te, t1...)
		t = append(t, te...)
	}
	return t
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

func compileAuthConfig(c *AuthConfig) []byte {
	var b []byte = []byte{0x04}
	b = append(b, bcd.FromUint(uint64(c.Amount), 6)...)
	if c.Currency != nil {
		b = append(b, 0x49)
		b = append(b, bcd.FromUint(uint64(*c.Currency), 2)...)
	}
	if c.PaymentType != nil {
		b = append(b, 0x19, *c.PaymentType)
	}
	if c.CardNumber != nil {
		b = append(b, 0x22)
		b = append(b, *c.CardNumber...)
	}
	if c.ExpiryDate != nil {
		b = append(b, 0x0E)
		b = append(b, (*c.ExpiryDate).getBCD()...)
	}
	if c.TLV != nil {
		b = append(b, c.TLV.Marshal()...)
	}
	return b
}
