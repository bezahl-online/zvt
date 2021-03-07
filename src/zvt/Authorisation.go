package zvt

import (
	"bezahl.online/zvt/src/apdu"
	"bezahl.online/zvt/src/apdu/bmp"
	"bezahl.online/zvt/src/instr"
	"github.com/albenik/bcd"
)

// Authorisation implents 06 01
// initiates a payment process
func (p *PT) Authorisation(config *AuthConfig) (*Response, error) {
	ctrlField := instr.Map["Authorisation"]
	return p.send(Command{ctrlField, compileAuthConfig(config)})
}

func compileAuthConfig(c *AuthConfig) apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{ // FIXME:
			{ID: 0x49, Data: bcd.FromUint16(uint16(*c.Currency))},
			{ID: 0x19, Data: []byte{*c.PaymentType}},
		},
	}
}
