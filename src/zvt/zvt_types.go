package zvt

import (
	"net"
	"sync"

	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/zvt/tlv"
	"github.com/albenik/bcd"
)

// BMP structure
type BMP struct {
	bmp  byte
	data []byte
}

// PT is the class
type PT struct {
	lock *sync.RWMutex
	conn net.Conn
}

// EUR currency code
const EUR = 978

// Config is the config struct
type Config struct {
	pwd          [3]byte
	config       byte
	currency     int // default EUR
	service      byte
	tlvContainer *tlv.Container
}

// CompileConfig return a compiled byte array of the configuration
func (c *Config) CompileConfig() apdu.DataUnit {
	var dataUnit apdu.DataUnit = apdu.DataUnit{}
	var b []byte = []byte{}
	b = append(b, c.pwd[0], c.pwd[1], c.pwd[2])
	b = append(b, byte(c.config))
	b = append(b, bcd.FromUint16(uint16(c.currency))...)
	b = append(b, 0x03, byte(c.service))
	dataUnit.Data = b
	dataUnit.TLVContainer = *c.tlvContainer
	return dataUnit
}

// type CardData struct {
// 	ExpiryDate     ExpiryDate
// 	SequenceNumber SequenceNumber
// 	EFID           EFID
// 	CardType       CardType
// }

// ExpiryDate structur
type ExpiryDate struct {
	Month int
	Year  int
}

func (e *ExpiryDate) getBCD() []byte {
	var b []byte = bcd.FromUint(uint64(e.Month), 1)
	b = append(b, bcd.FromUint(uint64(e.Year), 1)...)
	return b
}
