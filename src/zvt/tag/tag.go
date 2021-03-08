package tag

import (
	"fmt"

	"bezahl.online/zvt/src/apdu/bmp/blen"
)

// TODO: type InsideTAG []byte

// Info is the TAG info structure
type Info struct {
	Name       string
	LengthType int
	Length     int
	TAGNrLen   int
}

// InfoMaps is the TAG info maps collection
var InfoMaps *IMaps = &IMaps{
	InfoMap:  make(map[byte]Info),
	InfoMapE: make(map[[2]byte]Info),
}

func init() {
	InfoMaps.initInfoMap()
	InfoMaps.initInfoMapE()
}

func (m *IMaps) initInfoMap() {
	m.InfoMap[0x22] = Info{"PAN / EF_ID", blen.LL, 0, 1}
	m.InfoMap[0x24] = Info{"text message", blen.BINARY, 1, 1}
	m.InfoMap[0x25] = Info{"print-texts", blen.BINARY, 0, 1}
	m.InfoMap[0x26] = Info{"List  of  permitted  ZVT-commands", blen.BINARY, 0, 1}
	m.InfoMap[0x43] = Info{"application-ID (RID+PIX)", blen.BINARY, 0, 1}
	m.InfoMap[0x45] = Info{"receipt-parameter", blen.BINARY, 4, 1}
	m.InfoMap[0x46] = Info{"EMV-print-data (customer-receipt)", blen.BINARY, 0, 1}
	m.InfoMap[0x47] = Info{"EMV-print-data (merchant-receipt)", blen.BINARY, 0, 1}
	m.InfoMap[0x6F] = Info{"incorrect currency", blen.NONE, 0, 1}
}

func (m *IMaps) initInfoMapE() {
	// 0x01 transaction-receipt (merchant-receipt)
	// 0x02 transaction-receipt (customer-receipt)
	// 0x03 administration-receipt
	m.InfoMapE[[2]byte{0x1F, 0x04}] = Info{"receipt-parameter", blen.BINARY, 1, 2}
	m.InfoMapE[[2]byte{0x1F, 0x07}] = Info{"receipt-type", blen.BINARY, 1, 2}
	m.InfoMapE[[2]byte{0x1F, 0x10}] = Info{"cardholder authentication", blen.BINARY, 1, 2}
	m.InfoMapE[[2]byte{0x1F, 0x12}] = Info{"card-technology", blen.BINARY, 1, 2}
	m.InfoMapE[[2]byte{0x1F, 0x5B}] = Info{"timeout", blen.BINARY, 1, 2}
}

type IMaps struct {
	InfoMap  map[byte]Info
	InfoMapE map[[2]byte]Info
}

// GetInfoMap returns the Info depending on the
// first one or two bytes of the given data
func (m *IMaps) GetInfoMap(nr []byte) (Info, bool) {
	var info Info
	var found bool
	if nr[0]&0x1F == 0x1F {
		info, found = m.InfoMapE[[2]byte{nr[0], nr[1]}]
	} else {
		info, found = m.InfoMap[nr[0]]
	}
	return info, found
}

// Decompile retrieves the TAG number and data from the
// TLV data objects according to the ZTV protocol
func Decompile(data *[]byte) ([]byte, error) {
	d := *data
	if d[0]&0x1F == 0x1F {
		// in theory it could by another byte long
		// but it never happens
		if len(d) < 2 {
			return d[:1], fmt.Errorf("wrong TAG format: second byte expected")
		}
		return d[:2], nil
	}
	return d[:1], nil
}
