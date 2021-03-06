package tlv

import "fmt"

// DecompileTAG retrieves the TAG number and data from the
// TLV data objects according to the ZTV protocol
func DecompileTAG(data *[]byte) ([]byte, error) {
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
