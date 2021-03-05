package zvt

import (
	"net"
	"sync"

	"github.com/albenik/bcd"
)

// Command is the structur for a APDU
type Command struct {
	Class  byte
	Inst   byte
	Length byte
	Data   []byte
}

// BMP structure
type BMP struct {
	bmp  byte
	data []byte
}

// Response is the response from the PT
type Response struct {
	CCRC   byte
	APRC   byte
	Length int
	Data   []byte
}

// Marshal returns the bytes of the Response structure
func (r *Response) Marshal() []byte {
	var b []byte = []byte{r.CCRC, r.APRC, byte(r.Length)}
	b = append(b, r.Data...)
	return b
}

// IsACK returns true if the response is in fact a ACK
func (r *Response) IsACK() bool {
	return r.CCRC == 0x80 || (r.CCRC == 0x84 && r.APRC == 0)
}

// IsIntermediate returns true if the response is in fact a ACK
func (r *Response) IsIntermediate() bool {
	return r.CCRC == 0x04 && r.APRC == 0xFF
}

// IsStatus returns true if the response has a status byte
func (r *Response) IsStatus() bool {
	return r.CCRC == 0x04 && r.APRC == 0x0F
}

// GetTLV return a TLV or an error
func (r *Response) GetTLV() (*TLV, error) {
	var tlv TLV = TLV{
		Objects: []DataObject{},
	}
	err := tlv.Unmarshal(&r.Data)
	return &tlv, err
}

// PT is the class
type PT struct {
	lock *sync.RWMutex
	conn net.Conn
}

// TAG is the TAG structure
type TAG struct {
	Data []byte
}

// DataObject is part of a TLV
type DataObject struct {
	TAG  []byte
	data []byte
}

// TLV is the type length value container
type TLV struct {
	Objects []DataObject
}

const (
	// ServiceMenuNOTAssignedToFunctionKey prevents PT from assigning the service menu to the function key
	ServiceMenuNOTAssignedToFunctionKey = 1 << iota

	// DisplayTextsForCommandsAuthorisation Pre-initialisation and Reversal will be displayed in capitals
	DisplayTextsForCommandsAuthorisation
)

// EUR currency code
const EUR = 978

// Config is the config struct
type Config struct {
	pwd      [3]byte
	config   byte
	currency int // default EUR
	service  byte
	tlv      *TLV
}

// CompileConfig return a compiled byte array of the configuration
func (c *Config) CompileConfig() []byte {
	var b []byte = []byte{}
	b = append(b, c.pwd[0], c.pwd[1], c.pwd[2])
	b = append(b, byte(c.config))
	b = append(b, bcd.FromUint16(uint16(c.currency))...)
	b = append(b, 0x03, byte(c.service))
	if c.tlv != nil {
		b = append(b, c.tlv.Marshal()...)
	}
	return b
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

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int
	Currency    *int
	PaymentType *byte
	ExpiryDate  *ExpiryDate
	CardNumber  *[]byte
	TLV         *TLV
}

func (a *AuthConfig) marshal() {
	var b []byte = []byte{0x04}
	b = append(b, bcd.FromUint(uint64(a.Amount), 6)...)
	if a.Currency != nil {
		b = append(b, byte(0x49))
		b = append(b, bcd.FromUint(uint64(*a.Currency), 2)...)
	}
	if a.PaymentType != nil {
		b = append(b, 0x19, *a.PaymentType)
	}
	if a.ExpiryDate != nil {
		b = append(b, 0x0e)
		b = append(b, (*a.ExpiryDate).getBCD()...)
	}
}
