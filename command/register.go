package command

import (
	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
)

// Config is the config struct
type Config struct {
	pwd          [3]byte
	config       byte
	currency     int // default EUR
	service      byte
	tlvContainer *tlv.Container
}

// Register implements inst 06 00
// set up different configurations on the PT
func (p *PT) Register(config *Config) error {
	i := instr.Map["Registration"]
	return p.send(Command{
		CtrlField: i,
		Data:      (*config).CompileConfig(),
	})
}

// CompileConfig return a compiled byte array of the configuration
func (c *Config) CompileConfig() apdu.DataUnit {
	var dataUnit apdu.DataUnit = apdu.DataUnit{}
	var b []byte = []byte{}
	b = append(b, c.pwd[0], c.pwd[1], c.pwd[2])
	b = append(b, byte(c.config))
	b = append(b, bcd.FromUint16(uint16(c.currency))...)
	b = append(b, 0x03, byte(c.service))
	dataUnit.Data = b
	dataUnit.TLVContainer = *c.tlvContainer
	return dataUnit
}
