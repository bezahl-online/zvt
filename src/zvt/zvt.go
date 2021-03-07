package zvt

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/tlv"
	"bezahl.online/zvt/src/zvt/util"
)

// ZVT represents the driver
var ZVT PT

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

// Command is the structur for a APDU
type Command struct {
	Instr instr.CtrlField
	Data  apdu.DataUnit
}

// Marshal marshal every thing to final command
func (c *Command) Marshal() ([]byte, error) {
	var b []byte = []byte{}
	data, err := c.Data.Marshal()
	if err != nil {
		return b, err
	}
	b = append(b, c.Instr.Marshal(uint16(len(data)))...)
	b = append(b, data...)
	return b, nil
}

// Unmarshal is
func (c *Command) Unmarshal(data *[]byte) error {
	i := instr.Find(data)
	if i == nil {
		return fmt.Errorf("APRC % X not found", (*data)[:2])
	}
	c.Instr = *i
	dstart := 3
	if (*data)[3] == 0xff {
		dstart = 5
	}
	dend := dstart + i.RawDataLength
	raw := (*data)[dstart:dend]
	objs := []bmp.OBJ{}
	for {
		if len(*data) <= dend || (*data)[dend] == bmp.TLV {
			break
		}
		o := bmp.OBJ{}
		err := o.Unmarshal(((*data)[dend:]))
		if err != nil {
			return err
		}
		objs = append(objs, o)
		dend += o.Size
	}
	tlv := tlv.Container{}
	tlvData := (*data)[dend:]
	tlv.Unmarshal(&tlvData)
	c.Data = apdu.DataUnit{
		Data:         raw,
		BMPOBJs:      objs,
		TLVContainer: tlv,
	}
	return nil
}

// SendACK send ACK and return the response or error
func (p *PT) SendACK(timeout time.Duration) error {
	i := instr.Map["ACK"]
	_, err := p.send(Command{
		Instr: i,
	})
	if err != nil {
		return err
	}
	return nil
}

// FIXME: create comm package!

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

func (p *PT) send(c Command) (*Command, error) {
	var err error
	b, err := c.Marshal()
	if err != nil {
		return nil, err
	}
	fmt.Printf("% X", b)
	_, err = p.conn.Write(b)
	if err != nil {
		return nil, err
	}

	var resp *Command
	resp, err = p.readResponse(5 * time.Second)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (p *PT) readResponse(timeout time.Duration) (*Command, error) {
	var resp *Command
	var err error
	var readBuf []byte = make([]byte, 1024)
	p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err := p.conn.Read(readBuf)
	util.Save(&readBuf, nr)
	if err != nil {
		return resp, err
	}
	fmt.Printf("PT => ECR (%3d):% X\n", nr, readBuf[:nr])
	// resp, err = resp.Unmarshal(readBuf[:nr])
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
