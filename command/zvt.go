package command

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/util"
	"go.uber.org/zap"
)

var dontstartTest = false

func skipShort(t *testing.T) bool {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	if dontstartTest {
		fmt.Println("skipping test")
		return true
	}
	return false
}

// PaymentTerminal represents the driver
var PaymentTerminal PT = PT{
	RWMutex: sync.RWMutex{},
}

// EUR currency code
const EUR = 978

var Logger *zap.Logger

// PT is the class
type PT struct {
	sync.RWMutex
	Logger *zap.Logger
	conn   net.Conn
}

// stanard timeout for read from and write to PT
const defaultTimeout = 5 * time.Second

func init() {
	PaymentTerminal.Logger = getLogger()
	Logger = PaymentTerminal.Logger
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

// Close closes the connection to the PT
func (p *PT) ReConnect() error {
	defer p.Unlock()
	p.Lock()
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return err
		}
		p.conn = nil
	}
	p.tryConnect()
	return nil
}

func (p *PT) reconnectIfLost() error {
	if p.conn == nil {
		// seems to be connected again
		//p.SendACK()
		err := p.tryConnect()
		if err != nil {
			err = fmt.Errorf("lost connection to PT")
			time.Sleep(defaultTimeout)
		}
		return err
	}
	return nil
}

func (p *PT) tryConnect() error {
	go p.Connect()
	time.Sleep(100 * time.Millisecond)
	if p.conn != nil {
		return nil
	}
	return fmt.Errorf("not connected yet")
}

// SendACK send ACK and return the response or error
func (p *PT) SendACK() error {
	defer p.Unlock()
	p.Lock()
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

func (p *PT) SendCommand(command Command) error {
	defer p.Unlock()
	p.Lock()
	if err := p.send(command); err != nil {
		return p.logSendError(err)
	}
	response, err := PaymentTerminal.ReadResponse()
	if err != nil {
		return p.logResponseError(err)
	}
	return response.IsAck()
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
		if errors.Is(err, net.ErrClosed) {
			p.conn = nil // we lost connection - lets reconnect
			if err = p.reconnectIfLost(); err != nil {
				p.Logger.Error(err.Error())
				_, err = p.conn.Write(b)
				if err != nil {
					return err
				}
			}
		}
	}
	data, err := c.Data.Marshal()
	if err != nil {
		Logger.Error(err.Error())
	}
	util.Save(&data, &c.CtrlField, "EC")
	return nil
}

func (p *PT) logSendError(err error) error {
	p.Logger.Error("error while sending command to PT",
		zap.Error(err))
	return err
}

func (p *PT) logResponseError(err error) error {
	p.Logger.Error("error while reading response from PT",
		zap.Error(err))
	return err
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
		if errors.Is(err, net.ErrClosed) {
			p.conn = nil // we lost connection - lets reconnect
			if err = p.reconnectIfLost(); err != nil {
				p.Logger.Error(err.Error())
				n, err = p.conn.Read(cf)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	i := instr.Find(&cf) // FIXME: not tested from here to end
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
		message := fmt.Sprintf("PT => ECR [% 02X] (  0)", cf[:2])
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
