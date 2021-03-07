package tag

import (
	"bezahl.online/zvt/src/apdu/bmp/blen"
)

// Info is the TAG info structure
type Info struct {
	Name       string
	LengthType int
	Length     int
	TAGNrLen   int
}

// InfoMaps is the TAG info maps collection
var InfoMaps *IMaps = &IMaps{
	InfoMap:  make(map[byte]Info),
	InfoMapE: make(map[[2]byte]Info),
}

func init() {
	InfoMaps.initInfoMap()
	InfoMaps.initInfoMapE()
}

func (m *IMaps) initInfoMap() {
	m.InfoMap[0x24] = Info{"text message", blen.BINARY, 1, 1}
	m.InfoMap[0x6F] = Info{"incorrect currency", blen.NONE, 0, 1}
}

func (m *IMaps) initInfoMapE() {
	m.InfoMapE[[2]byte{0x1F, 0x5B}] = Info{"timeout", blen.BINARY, 1, 2}
}

type IMaps struct {
	InfoMap  map[byte]Info
	InfoMapE map[[2]byte]Info
}

// GetInfoMap returns the Info depending on the
// first one or two bytes of the given data
func (m *IMaps) GetInfoMap(nr []byte) (Info, bool) {
	var info Info
	var found bool
	if nr[0]&0x1F == 0x1F {
		info, found = m.InfoMapE[[2]byte{nr[0], nr[1]}]
	} else {
		info, found = m.InfoMap[nr[0]]
	}
	return info, found
}
