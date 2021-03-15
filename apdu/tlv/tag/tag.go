package tag

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu/bmp/blen"
)

// TODO: type InsideTAG []byte
// could map all the ways to interpret the
// BMP's data correctly

// Info is the TAG info structure
type Info struct {
	Name       string
	LengthType int
	Length     int
	TAGNrLen   int
}

// InfoMap holds all used commands
var InfoMap map[[2]byte]Info = make(map[[2]byte]Info)

func init() {
	InfoMap[[2]byte{0x22}] = Info{"PAN / EF_ID", blen.LL, 0, 1}
	InfoMap[[2]byte{0x24}] = Info{"text message", blen.BINARY, 1, 1}
	InfoMap[[2]byte{0x25}] = Info{"print-texts", blen.BINARY, 0, 1}
	InfoMap[[2]byte{0x26}] = Info{"List  of  permitted  ZVT-commands", blen.BINARY, 0, 1}
	InfoMap[[2]byte{0x43}] = Info{"application-ID (RID+PIX)", blen.BINARY, 0, 1}
	InfoMap[[2]byte{0x45}] = Info{"receipt-parameter", blen.BINARY, 4, 1}
	InfoMap[[2]byte{0x46}] = Info{"EMV-print-data (customer-receipt)", blen.BINARY, 0, 1}
	InfoMap[[2]byte{0x47}] = Info{"EMV-print-data (merchant-receipt)", blen.BINARY, 0, 1}
	InfoMap[[2]byte{0x6F}] = Info{"incorrect currency", blen.NONE, 0, 1}
	// 0x01 transaction-receipt (merchant-receipt)
	// 0x02 transaction-receipt (customer-receipt)
	// 0x03 administration-receipt
	InfoMap[[2]byte{0x1F, 0x04}] = Info{"receipt-parameter", blen.BINARY, 1, 2}
	InfoMap[[2]byte{0x1F, 0x07}] = Info{"receipt-type", blen.BINARY, 1, 2}
	InfoMap[[2]byte{0x1F, 0x10}] = Info{"cardholder authentication", blen.BINARY, 1, 2}
	InfoMap[[2]byte{0x1F, 0x12}] = Info{"card-technology", blen.BINARY, 1, 2}
	InfoMap[[2]byte{0x1F, 0x5B}] = Info{"timeout", blen.BINARY, 1, 2}
}

// Decompile retrieves the TAG number and data from the
// TLV data objects according to the ZTV protocol
func Decompile(data *[]byte) ([2]byte, error) {
	d := *data
	var tagNr [2]byte = [2]byte{}
	tagNr[0] = d[0]
	if d[0]&0x1F == 0x1F {
		// in theory it could by another byte long
		// but it never happens
		if len(d) < 2 {
			return tagNr, fmt.Errorf("wrong TAG format: second byte expected")
		}
		tagNr[1] = d[1] // add second byte
	}
	return tagNr, nil
}
