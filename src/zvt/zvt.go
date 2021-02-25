package zvt

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
func (c command) getBytes() []byte {
	var b []byte = []byte{
		c.Class,
		c.Inst,
		c.Length,
	}
	b[2] = byte(len(c.Data))
	b = append(b, c.Data...)
	return b
}

func (p *PT) send(c command) error {
	nr, err := p.conn.Write(c.getBytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("debug: %d bytes written:\n", nr)
	fmt.Println(strings.ReplaceAll(fmt.Sprintf("%q\n", fmt.Sprintf("% #x", c.getBytes())), " ", ","))

	var readBuf []byte = make([]byte, 3)
	p.conn.SetDeadline(time.Now().Add(5 * time.Second))
	nr, err = p.conn.Read(readBuf)
	if err != nil {
		log.Fatal(err)
	}
	if nr == 3 && readBuf[0] == 0x80 {
		p.conn.Write([]byte{0x80, 0, 0})
	} else {
		return fmt.Errorf("Error from PT: %x", readBuf)
	}
	return nil
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

func (p *PT) compileConfigByte(b []byte) ConfigByte {
	var cb ConfigByte = 0
	for _, v := range b {
		cb += ConfigByte(v)
	}
	return cb
}

func (p *PT) compileServiceByte(b []byte) ServiceByte {
	var sb ServiceByte = 0
	for _, v := range b {
		sb += ServiceByte(v)
	}
	return sb
}
