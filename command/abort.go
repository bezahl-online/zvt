package command

import (
	"github.com/bezahl-online/zvt/instr"
)

// Abort implements inst 06 B0
// ECR can instruct the PT to abort execution of a command
func (p *PT) Abort() error {
	Logger.Info("ABORT")
	return p.SendCommand(Command{CtrlField: instr.Map["Abort"]})
}
