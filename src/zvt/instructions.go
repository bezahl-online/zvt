package zvt

// DisplayText implements instr 06 E0
func (p *PT) DisplayText(text []string) error {
	return p.send(command{
		Class: 0x06,
		Inst:  0xe0,
		Data:  p.compileText(text),
	})
}

func (p *PT) Register(config *PTConfig) error {
	return p.send(command{
		Class: 0x06,
		Inst:  0x00,
		Data:  p.compileConfigBytes(config),
	})
}
