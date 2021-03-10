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

// PaymentTerminal represents the driver
var PaymentTerminal PT

// stanard timeout for read from and write to PT
const defaultTimeout = 5 * time.Second

func init() {
	var pt PT = PT{
		lock: &sync.RWMutex{},
		conn: nil,
	}
	err := pt.Open()
	if err != nil {
		fmt.Println(err.Error())
	}
	PaymentTerminal = pt
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
	if p.conn == nil {
		return fmt.Errorf("no connection to PT")
	}
	var err error
	b, err := c.Marshal()
	if err != nil {
		return err
	}
	fmt.Printf("ECR => PT (%3d):% X\n", len(b), b)
	p.conn.SetDeadline(time.Now().Add(defaultTimeout))
	nr, err := p.conn.Write(b)
	if err != nil {
		return err
	}
	if nr > 0 {
		util.Save(&b, nr, "EC")
	}
	return nil
}

// ReadResponse reads from the connection to the PT
func (p *PT) ReadResponse() (*Command, error) {
	return p.ReadResponseWithTimeout(defaultTimeout)
}

// ReadResponseWithTimeout reads from the connection to the PT
// where a timeout can be specified
// if reading time exceeds timout duration an error is returned
func (p *PT) ReadResponseWithTimeout(timeout time.Duration) (*Command, error) {
	var err error
	if p.conn == nil {
		return nil, fmt.Errorf("no connection to PT")
	}
	var resp *Command = &Command{}
	var cf []byte = []byte{0, 0, 0}
	p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err := p.conn.Read(cf)
	i := instr.Find(&cf)
	if i == nil {
		return nil, fmt.Errorf("control field '% X' not found", cf)
	}
	lenBuf := []byte{cf[2]}
	if cf[2] == 0xFF {
		lenBuf = append(lenBuf, 0, 0)
		nr, err = p.conn.Read(lenBuf[1:])
		if err != nil {
			return resp, err
		}
	}
	i.Length.Unmarshal(lenBuf)
	var readBuf []byte = make([]byte, i.Length.Value)
	// p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err = p.conn.Read(readBuf)
	if err != nil {
		return resp, err
	}
	fmt.Printf("PT => ECR (%03X):% X\n", []byte{i.Class, i.Instr}, readBuf[:nr])
	if nr > 0 {
		util.Save(&readBuf, nr, "PT")
	}
	data := []byte{i.Class, i.Instr}
	data = append(data, lenBuf...)
	data = append(data, readBuf[:nr]...)
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
