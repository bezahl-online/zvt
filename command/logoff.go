package command

import "github.com/bezahl-online/zvt/instr"

// LogOff implements inst 06 02
// with following consequences:
// the PT resets the Registrationconfig-byte to ‘86’
// and the PT may not send any more TLV-containers
func (p *PT) LogOff() error {
	return p.SendCommand(Command{
		CtrlField: instr.Map["LogOff"],
	})
}
