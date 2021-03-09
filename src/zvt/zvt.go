package zvt

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"bezahl.online/zvt/src/instr"
	"bezahl.online/zvt/src/zvt/util"
)

// EUR currency code
const EUR = 978

// PT is the class
type PT struct {
	lock *sync.RWMutex
	conn net.Conn
}

// ZVT represents the driver
var ZVT PT

// stanard timeout for read from and write to PT
const timeoutRW = 5 * time.Second

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

// SendACK send ACK and return the response or error
func (p *PT) SendACK() error {
	i := instr.Map["ACK"]
	err := p.send(Command{
		CtrlField: i,
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

func (p *PT) send(c Command) error {
	var err error
	b, err := c.Marshal()
	if err != nil {
		return err
	}
	fmt.Printf("ECR => PT (%3d):% X\n", len(b), b)
	p.conn.SetDeadline(time.Now().Add(timeoutRW))
	_, err = p.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// ReadResponse reads from the connection to the PT
func (p *PT) ReadResponse() (*Command, error) {
	var resp *Command = &Command{}
	var err error
	var readBuf []byte = make([]byte, 1024)
	p.conn.SetDeadline(time.Now().Add(timeoutRW))
	nr, err := p.conn.Read(readBuf)
	if nr > 0 {
		util.Save(&readBuf, nr)
	}
	if err != nil {
		return resp, err
	}
	fmt.Printf("PT => ECR (%3d):% X\n", nr, readBuf[:nr])
	data := readBuf[:nr]
	if nr == 1 && data[0] == 0x15 { // incorrect password
		return nil, fmt.Errorf("Incorrect Password")
	}
	err = resp.Unmarshal(&data)
	return resp, err
}

func compileLL(l uint8) []byte {
	var b []byte = make([]byte, 2)
	lz := uint8(l / 10)           // value of tens
	le := uint8(l - uint8(10*lz)) // value of unit position
	b[0] = 0xF0 + lz              // code into 0xFx (tens) (BCD)
	b[1] = 0xF0 + le              // code into 0xFy (unit) (BCD)
	return b
}
