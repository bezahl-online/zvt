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
type TLV struct {
	BMP  byte
	data []byte
}

// ServiceByte is the sructure of the ZVT service byte
type ServiceByte byte

const (
	// ServiceMenuNOTAssignedToFunctionKey prevents PT from assigning the service menu to the function key
	ServiceMenuNOTAssignedToFunctionKey = 1 << iota

	// DisplayTextsForCommandsAuthorisation Pre-initialisation and Reversal will be displayed in capitals
	DisplayTextsForCommandsAuthorisation
)

// ConfigByte is the structur of the ZVT config byte
type ConfigByte byte

const (
	// PaymentReceiptPrintedByECR ECR assumes receipt-printout for payment functions
	PaymentReceiptPrintedByECR = 2 << iota

	// AdminReceiptPrintedByECR ECR assumes receipt-printout for administration functions
	AdminReceiptPrintedByECR

	// PTSendsIntermediateStatus PTSendsIntermediateStatus
	PTSendsIntermediateStatus

	// AmountInputOnPTpossible ECR controls payment function
	AmountInputOnPTpossible

	// AdminFunctionOnPTpossible ECR controls administration function
	AdminFunctionOnPTpossible
	_

	// ECRusingPrintLinesForPrintout ECR print-type
	ECRusingPrintLinesForPrintout
)

// Currency the currency type
type Currency [2]byte

// EUR currency code
var EUR Currency = Currency([2]byte{0x09, 0x97})

// PTConfig ist the config struct
type PTConfig struct {
	pwd      [3]byte
	config   ConfigByte
	currency Currency
	service  ServiceByte
	tlv      *TLV
}
