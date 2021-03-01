package zvt

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/albenik/bcd"
)

// ZVT represents the driver
var ZVT PT
var zvtACK []byte = []byte{0x80, 0x00, 0x00}

func init() {
	var pt PT = PT{
		lock: &sync.RWMutex{},
		conn: nil,
	}

	err := pt.Open()
	if err != nil {
		log.Fatal(err)
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
func (c Command) getBytes() []byte {
	var b []byte = []byte{
		c.Class,
		c.Inst,
		c.Length,
	}
	b[2] = byte(len(c.Data))
	b = append(b, c.Data...)
	return b
}

func (p *PT) send(c Command) (Response, error) {
	nr, err := p.conn.Write(c.getBytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("debug: %d bytes written:\n", nr)
	fmt.Println(strings.ReplaceAll(fmt.Sprintf("%q\n", fmt.Sprintf("% #x", c.getBytes())), " ", ","))
	var resp Response
	var readBuf []byte = make([]byte, 128)
	p.conn.SetDeadline(time.Now().Add(5 * time.Second))
	nr, err = p.conn.Read(readBuf)
	if err != nil {
		return resp, err
	}
	resp, err = p.unmarshalAPDU(readBuf)
	p.conn.Write(zvtACK)
	return resp, err
}

func (p *PT) compileText(textarray []string) []byte {
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

func (p *PT) compilePTConfig(c *PTConfig) []byte {
	var b []byte = []byte{}
	b = append(b, c.pwd[0], c.pwd[1], c.pwd[2])
	b = append(b, byte(c.config))
	b = append(b, bcd.FromUint(uint64(c.currency), 2)...)
	b = append(b, 0x03, byte(c.service))
	if c.tlv != nil {
		b = append(b, p.marshalTLV(&c.tlv.Objects)...)
	}
	return b
}

func marshalDataObjects(dos *[]DataObject) []byte {
	var data []byte
	for _, obj := range *dos {
		data = append(data, obj.TAG...)
		data = append(data, compileLength(len(obj.data))...)
		data = append(data, obj.data...)
	}
	return data
}

const tlvBMP = 0x06

func (p *PT) marshalTLV(dos *[]DataObject) []byte {
	var b []byte
	b = append(b, tlvBMP)
	data := marshalDataObjects(dos)
	b = append(b, compileLength(len(data))...)
	b = append(b, data...)
	return b
}

func compileLength(len int) []byte {
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

func (p *PT) unmarshalAPDU(apduBytes []byte) (Response, error) {
	var resp Response
	if len(apduBytes) < 3 {
		return resp, fmt.Errorf("APDU less than 3 bytes long")
	}
	resp = Response{
		CCRC:   apduBytes[0],
		APRC:   apduBytes[1],
		Length: int(apduBytes[2]),
	}
	if len(apduBytes) >= int(apduBytes[2])+3 {
		resp.Data = apduBytes[3 : apduBytes[2]+3]
	}
	return resp, nil
}
