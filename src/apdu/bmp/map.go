package bmp

const (
	// Timeout in seconds
	// 1 byte binary
	Timeout = 0x01

	// MaxNrStatus information
	// 1 byte binary
	MaxNrStatus = 0x02

	// ServiceByte meaning depends ont command
	// 1 byte binary (bit field)
	ServiceByte = 0x03

	// Amount in cent
	// 6 byte BCD
	Amount = 0x04

	// PumpNumber is
	// 1 byte binary
	PumpNumber = 0x05

	// TLV encoded container
	TLV = 0x06

	// TraceNumer is
	// 3 byte BCD
	TraceNumer = 0x0B

	// Time HHMMSS
	// 	3 byte BCD
	Time = 0x0C

	// Date MMDD
	// 2 byte BCD
	Date = 0x0D

	// ExpiryDate YYMM
	// 2 byte BCD
	ExpiryDate = 0x0E

	// CardSeqNr number
	// 2 byte BCD
	CardSeqNr = 0x17

	// StatusByte (06 00)
	// 1 byte binary (bit field)
	StatusByte = 0x19

	// PaymentType (06 01)
	// 1 byte binary (bit field)
	PaymentType = 0x19

	// CardType (06 C0)
	// 1 byte binary (bit field)
	CardType = 0x19

	// PANEFID or EF_ID E' used to indicate a masked numeric digit1
	// If the card-number contains an odd number of digits, it is padded with an ‘F’
	// LL-Var BCD (Fx Fy) len = 10x+y
	PANEFID = 0x22

	// TrackData2 without start and end markers
	// 'E' used to indicate a masked numeric digit
	// LL-Var
	TrackData2 = 0x23
	// TrackData3 without start and end markers
	// 'E' used to indicate a masked numeric digit
	// LLL-Var
	TrackData3 = 0x24

	// ResultCode error code
	// 1 byte binary
	ResultCode = 0x27

	// TerminalID identifyer
	// 4 byte BCD
	TerminalID = 0x29

	// VUNr is
	// 15 byte ASCII
	VUNr = 0x2A

	// Currency EUR=978
	// 2 byte BCD
	Currency = 0x49
)

const (
	// NONE no legth (tag has no data)
	NONE = iota
	// BINARY legth binary coded
	BINARY
	// LL length 0xFx,0xFy -> BCD coded (10x+y)
	LL
	// LLL length 0xFx,0xFy,0xFz -> BCD coded (100x+10y+z)
	LLL
	// BCD fixed length depending on TAG BCD coded
	BCD
)

// Info is the TAG info structure
type Info struct {
	Name       string
	LengthType int
	FixLen     int
}

// InfoMap maps all used BMPs
var InfoMap map[byte]Info = make(map[byte]Info)

func init() {
	//                    Name      LenType      FixLen
	InfoMap[0x04] = Info{"amount in cent", NONE, 6}
	InfoMap[0x05] = Info{"pump number", NONE, 1}
	InfoMap[0xF1] = Info{"Text1 line 1", LL, 0}
	InfoMap[0xF2] = Info{"Text1 line 2", LL, 0}
	InfoMap[0xF3] = Info{"Text1 line 3", LL, 0}
	InfoMap[0xF4] = Info{"Text1 line 4", LL, 0}
}
