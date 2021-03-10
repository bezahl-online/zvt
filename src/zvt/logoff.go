package zvt

import "github.com/bezahl-online/zvt/src/instr"

// LogOff implements inst 06 02
// with following consequences:
// the PT resets the Registrationconfig-byte to ‘86’
// and the PT may not send any more TLV-containers
func (p *PT) LogOff() error {
	i := instr.Map["LogOff"]
	return p.send(Command{
		CtrlField: i,
	})
}
