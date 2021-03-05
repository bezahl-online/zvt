package zvt

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"bezahl.online/zvt/src/zvt/bmp"
	"bezahl.online/zvt/src/zvt/tag"
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

func marshalDataObjects(dos *[]DataObject) []byte {
	var data []byte
	for _, obj := range *dos {
		data = append(data, obj.TAG...)
		data = append(data, compileLength(len(obj.data))...)
		data = append(data, obj.data...)
	}
	return data
}

// Marshal retuns the byte array of the tlv
func (t *TLV) Marshal() []byte {
	var b []byte
	b = append(b, bmp.TLV)
	data := marshalDataObjects(&t.Objects)
	b = append(b, compileLength(len(data))...)
	b = append(b, data...)
	return b
}

// Unmarshal fills the structur with the given data
func (t *TLV) Unmarshal(data *[]byte) error {
	d := *data
	if d[0] == bmp.TLV {
		lenData := d[1:5]
		tlvLen, sizeOfLenField, err := decompileLength(&lenData)
		if err != nil {
			return err
		}
		if uint16(len(d))-sizeOfLenField == tlvLen {
			return fmt.Errorf("value in length field (%d) and length of data (%d) does not match", tlvLen, len(d))
		}
		// reduce data to data after TLV TAG (06) and length byte(s)
		d = d[sizeOfLenField+1:]
		if t.Objects == nil {
			t.Objects = []DataObject{}
		}
		for {
			obj, objLength, err := unmarshalDataObject(d)
			if err != nil {
				return err
			}
			t.Objects = append(t.Objects, obj)
			if len(d) == int(objLength) {
				break
			}
			d = d[objLength:]
		}
	}
	return nil
}

func unmarshalDataObject(d []byte) (DataObject, uint16, error) {
	tag, err := decompileTAG(&d)
	if err != nil {
		return DataObject{}, 0, err
	}
	tagLength := uint16(len(tag))
	tagLengthData := d[tagLength:]
	tagDataLength, tagLengthSize, err := decompileLength(&tagLengthData)
	objectLength := tagLength + tagLengthSize + tagDataLength
	d = d[tagLength+tagLengthSize:]
	obj := DataObject{
		TAG:  tag,
		data: d[:tagDataLength],
	}
	return obj, objectLength, nil
}

func decompileTAG(data *[]byte) ([]byte, error) {
	d := *data
	if d[0]&0x1F == 0x1F {
		// in theory it could by another byte long
		// but it never happens
		if len(d) < 2 {
			return d[:1], fmt.Errorf("wrong TAG format: second byte expected")
		}
		return d[:2], nil
	}
	return d[:1], nil
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

func decompileLength(data *[]byte) (uint16, uint16, error) {
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
	var resp Response
	if len(apduBytes) < 3 {
		return &resp, fmt.Errorf("APDU less than 3 bytes long")
	}
	resp = Response{
		CCRC:   apduBytes[0],
		APRC:   apduBytes[1],
		Length: int(apduBytes[2]),
	}
	if len(apduBytes) >= int(apduBytes[2])+3 {

		resp.Data = apduBytes[3 : apduBytes[2]+3]
	}
	if apduBytes[0] == 0x04 {
		if apduBytes[1] == 0xFF {
			resp.IStatus = apduBytes[3]
			if len(apduBytes) < 5 {
				return &resp, nil
			}
			d := resp.Data[4:]
			var dataStart, dataEnd, tagDataLength int
			tagNr := d[:2]
			tagInfo := tag.InfoMaps.GetInfoMap(tagNr)
			switch tagInfo.LengthType {
			case tag.BINARY:
				tagDataLength = int(d[tagInfo.TAGNrLen])
			}
			dataStart = tagInfo.TAGNrLen + tagInfo.Length
			dataEnd = dataStart + int(tagDataLength)
			resp.TLV = TLV{
				Objects: []DataObject{},
			}
			resp.TLV.Objects = append(resp.TLV.Objects, DataObject{
				TAG:  []byte{0x24},
				data: d[dataStart:dataEnd],
			})
		}
	}
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
