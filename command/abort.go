package command

import (
	"github.com/bezahl-online/zvt/instr"
)

// Abort implements inst 06 B0
// ECR can instruct the PT to abort execution of a command
func (p *PT) Abort() error {
	Logger.Info("ABORT")
	if err := p.send(Command{CtrlField: instr.Map["Abort"]}); err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err != nil {
		return err
	}
	return response.IsAck()
}
