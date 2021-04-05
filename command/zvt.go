package command

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"go.uber.org/zap"
)

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

// PaymentTerminal represents the driver
var PaymentTerminal PT = PT{
	Logger: &zap.Logger{},
	conn:   nil,
}

// EUR currency code
const EUR = 978

// PT is the class
type PT struct {
	Logger *zap.Logger
	conn   net.Conn
}

// stanard timeout for read from and write to PT
const defaultTimeout = 5 * time.Second

func init() {
	initLogger()
	// Logger.Debug("logger initialized")
	// var pt PT = PT{
	// 	lock: &sync.RWMutex{},
	// 	conn: nil,
	// }
	// PaymentTerminal = pt
}

// SendACK send ACK and return the response or error
func (p *PT) SendACK() error {
	p.Logger.Info("ACK")
	i := instr.Map["ACK"]
	err := p.send(Command{
		CtrlField: i,
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: create comm interface!

// Open opens a connection to the PT
func (p *PT) Open() error {
	var url string
	if len(os.Getenv("ZVT_URL")) > 0 {
		url = os.Getenv("ZVT_URL")
	} else {
		return fmt.Errorf("please set environment variabel ZVT_URL")
	}
	var err error
	p.conn, err = net.Dial("tcp", url)
	if err != nil {
		return err
	}
	p.conn.SetDeadline(time.Now().Add(5 * time.Second))
	return nil
}

func (p *PT) reconnectIfLost() error {
	if p.conn == nil {
		go p.Connect()
		time.Sleep(10 * time.Millisecond)
		if p.conn != nil {
			// seems to be connected again
			return nil
		}
		err := fmt.Errorf("lost connection to PT")
		return err
	}
	return nil
}

func (p *PT) send(c Command) error {
	log.Println("will write")
	var err error
	if err = p.reconnectIfLost(); err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	b, err := c.Marshal()
	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	logCommand(c, b)
	p.conn.SetDeadline(time.Now().Add(defaultTimeout))
	_, err = p.conn.Write(b)
	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	util.Save(&[]byte{}, &c.CtrlField, "EC")
	return nil
}

func logCommand(c Command, b []byte) {
	cf := c.CtrlField
	l := len(b)
	more := ""
	if l > 20 {
		l = 20
		more = "..."
	}
	Logger.Debug(fmt.Sprintf("ECR => PT [% 02X] Data: % 02X%s", []byte{cf.Class, cf.Instr}, b[:l], more),
		zap.Int("len", len(b)),
		zap.String("data", byteArrayToHexString(b)),
	)
}

func byteArrayToHexString(b []byte) string {
	return fmt.Sprintf("% 02X", b)
}

// ReadResponse reads from the connection to the PT
func (p *PT) ReadResponse() (*Command, error) {
	return p.ReadResponseWithTimeout(defaultTimeout)
}

// ReadResponseWithTimeout reads from the connection to the PT
// where a timeout can be specified
// if reading time exceeds timout duration an error is returned
func (p *PT) ReadResponseWithTimeout(timeout time.Duration) (*Command, error) {
	log.Println("will read")
	var err error
	if err = p.reconnectIfLost(); err != nil {
		p.Logger.Error(err.Error())
		return nil, err
	}
	var resp *Command = &Command{}
	var cf []byte = []byte{0, 0, 0}
	p.conn.SetDeadline(time.Now().Add(timeout))
	nr, err := p.conn.Read(cf)
	if err != nil {
		p.Logger.Error(err.Error())
		return resp, err
	}
	i := instr.Find(&cf)
	if i == nil {
		err := fmt.Errorf("control field '% X' not found", cf)
		p.Logger.Error(err.Error())
		return nil, err
	}
	lenBuf := []byte{cf[2]}
	if cf[2] == 0xFF {
		lenBuf = append(lenBuf, 0, 0)
		nr, err = p.conn.Read(lenBuf[1:])
		if err != nil {
			p.Logger.Error(err.Error())
			return resp, err
		}
	}
	i.Length.Unmarshal(lenBuf)
	if i.Length.Value == 0 {
		p.Logger.Debug("PT => ECR",
			zap.Int("len", nr),
			zap.String("data", byteArrayToHexString(cf)),
		)
		return &Command{
			CtrlField: *i,
		}, err
	}
	var readBuf []byte = make([]byte, i.Length.Value)
	nr, err = p.conn.Read(readBuf)
	if err != nil {
		p.Logger.Error(err.Error())
		return resp, err
	}
	data := []byte{i.Class, i.Instr}
	data = append(data, lenBuf...)
	data = append(data, readBuf[:nr]...)
	err = resp.Unmarshal(&data)
	util.Save(&data, i, "PT")
	p.Logger.Debug("PT => ECR",
		zap.Int("len", nr),
		zap.String("data", byteArrayToHexString(data)),
	)
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
