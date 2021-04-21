package command

import (
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/instr"
)

// Status implements inst 05 01
// ECR can request the Status of the PT
func (p *PT) Status() error {
	Logger.Info("STATUS")
	d := []byte(fixedPassword[:])
	d = append(d, 0x03, 0x07) // send TLV
	if err := p.send(Command{CtrlField: instr.Map["Status"],
		Data: apdu.DataUnit{
			Data: d,
		},
	}); err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err != nil {
		return err
	}
	return response.IsAck()
}
