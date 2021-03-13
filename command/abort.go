package command

import (
	"fmt"

	"github.com/bezahl-online/zvt/instr"
)

var fixedPassword [3]byte = [3]byte{0x12, 0x34, 0x56}

// Abort implements inst 06 B0
// ECR can instruct the PT to abort execution of a command
func (p *PT) Abort() error {
	if err := p.send(Command{CtrlField: instr.Map["Abort"]}); err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err == nil && !response.IsAck() {
		err = fmt.Errorf("error code %0X %0X", response.CtrlField.Class, response.CtrlField.Instr)
	}
	return nil
}
