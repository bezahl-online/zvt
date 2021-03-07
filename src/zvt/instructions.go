package zvt

var fixedPassword [3]byte = [3]byte{0x12, 0x34, 0x56}

// DisplayText implements instr 06 E0
func (p *PT) DisplayText(text []string) (*Command, error) {
	return p.send(Command{
		// Class: 0x06,
		// Inst:  0xe0,
		// Data:  compileText(text),
	})
}

// Register implements inst 06 00
// set up different configurations on the PT
func (p *PT) Register(config *Config) (*Command, error) {
	return p.send(Command{
		// Class: 0x06,
		// Inst:  0x00,
		Data: (*config).CompileConfig(),
	})
}

// Abort implements inst 06 B0
// ECR can instruct the PT to abort execution of a command
func (p *PT) Abort() (*Command, error) {
	return p.send(Command{
		// Class: 0x06,
		// Inst:  0xB0,
	})
}

// // FIXME: getting error "0x83 function not possible" from PT
// // ChangePassword (06 95)
// // change the merchant password required for some ZVT commands to the PT
// func (p *PT) ChangePassword(old, new *[3]byte) (*Response, error) {
// 	var data []byte = []byte{}
// 	data = append(data, (*old)[0:3]...)
// 	data = append(data, (*new)[0:3]...)
// 	return p.send(Command{
// 		// Class: 0x06,
// 		// Inst:  0x95,
// 		// Data:  data,
// 	})
// }

// LogOff implements inst 06 02
// with following consequences:
// the PT resets the Registrationconfig-byte to ‘86’
// and the PT may not send any more TLV-containers
func (p *PT) LogOff() (*Command, error) {
	return p.send(Command{
		// Class: 0x06,
		// Inst:  0x02,
	})
}
