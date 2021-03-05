package tag

// Info is the TAG info structure
type Info struct {
	Name       string
	LengthType int
	Length     int
	TAGNrLen   int
}

const (
	// NONE no legth (tag has no data)
	NONE = iota
	// BINARY legth binary coded
	BINARY
	// LL length 0xFx,0xFy -> BCD coded (10x+y)
	LL
	// LLL length 0xFx,0xFy,0xFz -> BCD coded (100x+10y+z)
	LLL
	// BCD fixed length depending on TAG BCD coded
	BCD
)

// InfoMap is the TAG info map
var InfoMaps *IMaps = &IMaps{
	InfoMap:  make(map[byte]Info),
	InfoMapE: make(map[[2]byte]Info),
}

func init() {
	InfoMaps.initInfoMap()
	InfoMaps.initInfoMapE()
}

func (m *IMaps) initInfoMap() {
	m.InfoMap[0x24] = Info{"text message", BINARY, 1, 1}
	m.InfoMap[0x6F] = Info{"incorrect currency", NONE, 0, 1}
}

func (m *IMaps) initInfoMapE() {
	m.InfoMapE[[2]byte{0x1F, 0x5B}] = Info{"timeout", BINARY, 1, 2}
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
