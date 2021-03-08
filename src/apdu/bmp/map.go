package bmp

import (
	"bezahl.online/zvt/src/apdu/bmp/blen"
)

// Info is the TAG info structure
type Info struct {
	Name   string
	Length blen.Length
}

// InfoMap maps all used BMPs
var InfoMap map[byte]Info = make(map[byte]Info)

func init() {
	//                    Name      LenType      FixLen
	InfoMap[0x04] = Info{"amount in cent", blen.Length{Kind: blen.NONE, Value: 6}}
	InfoMap[0x05] = Info{"pump number", blen.Length{Kind: blen.NONE, Value: 1}}
	// InfoMap[0x06] = Info{"TLV", blen.Length{Kind: blen.BINARY, Value: 1}} // dont MAP TLV is a PSEUDO BITMAP!
	InfoMap[0x0E] = Info{"Exp.Date", blen.Length{Kind: blen.NONE, Size: 2}}
	InfoMap[0xF1] = Info{"Text1 line 1", blen.Length{Kind: blen.LL}}
	InfoMap[0xF2] = Info{"Text1 line 2", blen.Length{Kind: blen.LL}}
	InfoMap[0xF3] = Info{"Text1 line 3", blen.Length{Kind: blen.LL}}
	InfoMap[0xF4] = Info{"Text1 line 4", blen.Length{Kind: blen.LL}}
}
