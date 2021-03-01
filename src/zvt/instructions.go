package zvt

var fixedPassword [3]byte = [3]byte{0x12, 0x34, 0x56}

// DisplayText implements instr 06 E0
func (p *PT) DisplayText(text []string) (Response, error) {
	return p.send(Command{
		Class: 0x06,
		Inst:  0xe0,
		Data:  p.compileText(text),
	})
}

// Register implements inst 06 00
// set up different configurations on the PT
func (p *PT) Register(config *PTConfig) (Response, error) {
	return p.send(Command{
		Class: 0x06,
		Inst:  0x00,
		Data:  p.compilePTConfig(config),
	})
}

// Authorisation implents 06 01
// initiates a payment process
// func (p *PT) Authorisation(data *AuthData) (Response, error) {
// 	return p.send(Command{
// 		Class: 0x06,
// 		Inst:  0x01,
// 		Data:  p.marshalAuthData(data),
// 	})
// }
