package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// URL is the default URL
const URL = "192.168.30.112:20007"

type command struct {
	Class  byte
	Inst   byte
	Length byte
	Data   []byte
}

// PT is the class
type PT struct {
	lock sync.RWMutex
	conn net.Conn
}

func (p *PT) Open() error {
	url := URL
	if len(os.Getenv("ZVT_URL")) > 3 {
		url = os.Getenv("ZVT_URL")
	}
	var err error
	p.conn, err = net.Dial("tcp", url)
	p.conn.SetDeadline(time.Now().Add(5 * time.Second))
	return err
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

var errNotConnected error = fmt.Errorf("zvt device not connected")

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

func main() {
	var pt PT = PT{
		lock: sync.RWMutex{},
		conn: nil,
	}

	err := pt.Open()
	if err != nil {
		fmt.Printf("\n*** Error while connection to gm65 scanner\n")
	}
	data := pt.compileText([]string{"Da steh ich nun,",
		"ich armer Tor,",
		"Und bin so klug",
		"als wie zuvor."})
	c := command{
		Class: 0x06,
		Inst:  0xe0,
		Data:  data,
	}
	pt.send(c)

}
