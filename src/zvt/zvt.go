package zvt

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

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
	// err := pt.Open()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
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
	nr, err := p.conn.Write(c.getBytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ECR => PT (%3d):% X\n", nr, c.getBytes())
	// fmt.Println(strings.ReplaceAll(fmt.Sprintf("%q\n", fmt.Sprintf("% #x", c.getBytes())), " ", ","))
	var resp *Response
	resp, err = p.readResponse(5 * time.Second)
	if err != nil {
		return resp, err
	}
	if reflect.DeepEqual(*resp, zvtACK) {
		resp, err = p.SendACK(5 * time.Second)
	}
	return resp, err
}

func (p *PT) readResponse(timeout time.Duration) (*Response, error) {
	var resp *Response
	var err error
	var readBuf []byte = make([]byte, 128)
	p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err := p.conn.Read(readBuf)
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

// Marshal retuns the byte array of the tlv
func (t *TLV) Marshal() []byte {
	var b []byte
	b = append(b, tlvBMP)
	data := marshalDataObjects(&t.Objects)
	b = append(b, compileLength(len(data))...)
	b = append(b, data...)
	return b
}

// Unmarshal fills the structur with the given data
func (t *TLV) Unmarshal(data *[]byte) error {
	d := *data
	if d[0] == 0x06 {
		tlvLen, sizeOfLenField, err := decompileLength(data)
		if err != nil {
			return err
		}
		if len(d)-sizeOfLenField == int(tlvLen) {
			return fmt.Errorf("value in length field (%d) and length of data (%d) does not match", tlvLen, len(d))
		}

	}
	return nil
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

func decompileLength(data *[]byte) (uint16, int, error) {
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

func (p *PT) unmarshalAPDU(apduBytes []byte) (*Response, error) {
	var resp *Response
	if len(apduBytes) < 3 {
		return resp, fmt.Errorf("APDU less than 3 bytes long")
	}
	resp = &Response{
		CCRC:   apduBytes[0],
		APRC:   apduBytes[1],
		Length: int(apduBytes[2]),
	}
	if len(apduBytes) >= int(apduBytes[2])+3 {
		resp.Data = apduBytes[3 : apduBytes[2]+3]
	}
	return resp, nil
}

func compilePTConfig(c *PTConfig) []byte {
	var b []byte = []byte{}
	b = append(b, c.pwd[0], c.pwd[1], c.pwd[2])
	b = append(b, byte(c.config))
	b = append(b, bcd.FromUint(uint64(c.currency), 2)...)
	b = append(b, 0x03, byte(c.service))
	if c.tlv != nil {
		b = append(b, c.tlv.Marshal()...)
	}
	return b
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
