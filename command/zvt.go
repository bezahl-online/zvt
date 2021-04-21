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
var PaymentTerminal PT = PT{}

// EUR currency code
const EUR = 978

var Logger *zap.Logger

// PT is the class
type PT struct {
	Logger *zap.Logger
	conn   net.Conn
}

// stanard timeout for read from and write to PT
const defaultTimeout = 5 * time.Minute

func init() {
	PaymentTerminal.Logger = getLogger()
	Logger = PaymentTerminal.Logger
}

// SendACK send ACK and return the response or error
func (p *PT) SendACK() error {
	p.Logger.Info("ECR: 'ACK'")
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
	logCommand(false, c, b)
	p.conn.SetDeadline(time.Now().Add(defaultTimeout))
	_, err = p.conn.Write(b)
	if err != nil {
		p.Logger.Error(err.Error())
		return err
	}
	data, err := c.Data.Marshal()
	if err != nil {
		Logger.Error(err.Error())
	}
	util.Save(&data, &c.CtrlField, "EC")
	return nil
}

func logCommand(fromPT bool, c Command, b []byte) {
	cf := c.CtrlField
	blen := len(b)
	apduStart := blen - int(c.CtrlField.Length.Value)
	from := "ECR => PT"
	if fromPT {
		from = "PT => ECR"
	}
	message := fmt.Sprintf("%s [% 02X]", from, []byte{cf.Class, cf.Instr})
	logstr := fmt.Sprintf("%s (%3d)", message, int(c.CtrlField.Length.Value))
	if blen > 3 {
		logstr += fmt.Sprintf(" APDU: % 02X", b[apduStart:])
	}
	log.Println(logstr)
	if blen > apduStart {
		l := blen
		andMore := ""
		if blen > 20 {
			l = 20
			andMore = "..."
		}
		data := b[apduStart:l]
		message = fmt.Sprintf("%s APDU: % 02X%s", message, data, andMore)
	}
	Logger.Debug(message,
		zap.Int("len", int(c.CtrlField.Length.Value)),
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
	var err error
	if err = p.reconnectIfLost(); err != nil {
		p.Logger.Error(err.Error())
		return nil, err
	}
	var resp *Command = &Command{}
	var cf []byte = []byte{0, 0, 0}
	deadline := time.Now().Add(timeout)
	p.conn.SetDeadline(deadline)
	n, err := p.conn.Read(cf)
	_ = n
	if err != nil {
		// p.Logger.Error(err.Error())
		return resp, err
	}
	i := instr.Find(&cf)
	if i == nil {
		err := fmt.Errorf("control field '% X' not found", cf)
		// p.Logger.Error(err.Error())
		return nil, err
	}
	lenBuf := []byte{cf[2]}
	if cf[2] == 0xFF {
		lenBuf = append(lenBuf, 0, 0)
		_, err = p.conn.Read(lenBuf[1:])
		if err != nil {
			// p.Logger.Error(err.Error())
			return resp, err
		}
	}
	i.Length.Unmarshal(lenBuf)
	if i.Length.Value == 0 {
		message := fmt.Sprintf("PT => ECR [% 02X]", cf[:2])
		p.Logger.Debug(message)
		log.Print(message)
		return &Command{
			CtrlField: *i,
		}, err
	}
	var readBuf []byte = make([]byte, i.Length.Value)
	nr, err := p.conn.Read(readBuf)
	if err != nil {
		// p.Logger.Error(err.Error())
		return resp, err
	}
	util.Save(&readBuf, i, "PT")
	data := []byte{i.Class, i.Instr}
	data = append(data, lenBuf...)
	data = append(data, readBuf[:nr]...)
	err = resp.Unmarshal(&data)
	logCommand(true, Command{CtrlField: *i}, data)
	return resp, err
}
