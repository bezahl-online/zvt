package command

import (
	"fmt"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
)

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int64
	PaymentType byte
	Currency    int
	// TLV         *tlv.Container
}

// Authorisation implents 06 01
// initiates a payment process
func (p *PT) Authorisation(config *AuthConfig) error {
	ctrlField := instr.Map["Authorisation"]
	err := p.send(Command{ctrlField, config.marshal()})
	if err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err == nil && !response.IsAck() {
		err = fmt.Errorf("error code %0X %0X", response.CtrlField.Class, response.CtrlField.Instr)
	}
	return err
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
		},
	}
}
