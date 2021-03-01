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
const EUR = 997

// PTConfig is the config struct
type PTConfig struct {
	pwd      [3]byte
	config   byte
	currency int // default EUR
	service  byte
	tlv      *TLV
}

// PaymentType defines the mode of payment
type PaymentType byte

const (
	// UseDataFromPreviousReadCard The PT should execute the payment using the data from the previous „Read Card“ command.
	//If no card-data is available, the PT sets the corresponding return-code in the Status-Information.
	UseDataFromPreviousReadCard = 2

	// PrinterReady (mainly used for evaluation tests)
	PrinterReady = 4

	// TippableTransaction (since DCPOS 2.5: ignored for EMV tip/tippable transactions)
	TippableTransaction = 8

	// Geldkarte for GiroCard (ignored for DC POS realted or other cards)
	Geldkarte = 16

	// OnlineWithoutPIN (OLV or EuroELV, if only EuroELV is supported by PT)
	// (ignored by DC POS related or other cards)
	OnlineWithoutPIN = 32

	// GirocardTransaction according to TA7.0 rules for TA 7.0 capable PTs
	// DC POS transaction for capable PT's otherwise ignored or refused
	// PIN based transaction for other cards
	GirocardTransaction = 48

	// PaymentExcludeGeldKarte Payment according to PTs decision excluding GeldKarte
	PaymentExcludeGeldKarte = 64

	// PaymentIncludeGeldKarte Payment according to PTs decision including GeldKarte
	PaymentIncludeGeldKarte = 65
)

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
	var b []byte = bcd.FromUint(uint64(e.Month), 2)
	b = append(b, bcd.FromUint(uint64(e.Year), 2)...)
	return b
}

// AuthData is the auth data struct
type AuthData struct {
	Amount      int
	Currency    *int
	PaymentType *byte
	ExpiryDate  *ExpiryDate
	CardNumber  *int
	TLV         *TLV
}

func (a *AuthData) marshal() {
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
	b = append(b)
}
