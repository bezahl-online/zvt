package zvt

// Authorisation implents 06 01
// initiates a payment process
func (p *PT) Authorisation(config *AuthConfig) (*Response, error) {
	// TODO: listen to PT and send ack - need to communicate also with api client
	return p.send(Command{
		Class: 0x06,
		Inst:  0x01,
		Data:  compileAuthConfig(config),
	})
}
