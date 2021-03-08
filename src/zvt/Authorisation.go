package zvt

import (
	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"github.com/albenik/bcd"
)

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int
	PaymentType byte
	Currency    int
	// TLV         *tlv.Container
}

// Authorisation implents 06 01
// initiates a payment process
func (p *PT) Authorisation(config *AuthConfig) error {
	ctrlField := instr.Map["Authorisation"]
	return p.send(Command{ctrlField, config.marshal()})
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
			// {ID: 0x49, Data: bcd.FromUint16(uint16(a.Currency))},
			// {ID: 0x19, Data: []byte{a.PaymentType}},
		},
	}
}
