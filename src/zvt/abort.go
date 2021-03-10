package zvt

import (
	"github.com/bezahl-online/zvt/src/instr"
)

var fixedPassword [3]byte = [3]byte{0x12, 0x34, 0x56}

// Abort implements inst 06 B0
// ECR can instruct the PT to abort execution of a command
func (p *PT) Abort() error {
	i := instr.Map["Abort"]
	return p.send(Command{
		CtrlField: i,
	})
}
