package zvt

import (
	"net"
	"sync"
)

type command struct {
	Class  byte
	Inst   byte
	Length byte
	Data   []byte
}

// PT is the class
type PT struct {
	lock *sync.RWMutex
	conn net.Conn
}

// TLV is the type length value container
type TLV struct{}

// PTConfig ist the config struct
type PTConfig struct {
	config  byte
	service byte
	tlv     TLV
}
