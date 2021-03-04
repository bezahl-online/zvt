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

// PTConfig is the config struct
type PTConfig struct {
	pwd      [3]byte
	config   byte
	currency int // default EUR
	service  byte
	tlv      *TLV
}

// type CardData struct {
// 	ExpiryDate     ExpiryDate
// 	SequenceNumber SequenceNumber
// 	EFID           EFID
// 	CardType       CardType
// }

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
