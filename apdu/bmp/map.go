package bmp

import (
	"github.com/bezahl-online/zvt/apdu/bmp/blen"
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
	InfoMap[0x0B] = Info{"Trace number", blen.Length{Kind: blen.NONE, Value: 3}}
	InfoMap[0x0C] = Info{"Time HHMMSS", blen.Length{Kind: blen.NONE, Value: 3}}
	InfoMap[0x0D] = Info{"Date MMDD", blen.Length{Kind: blen.NONE, Value: 2}}
	InfoMap[0x0E] = Info{"Exp.Date", blen.Length{Kind: blen.NONE, Value: 2}}
	InfoMap[0x17] = Info{"Card sequence-number", blen.Length{Kind: blen.NONE, Value: 2}}
	InfoMap[0x19] = Info{"Status/PaymentType/CardType", blen.Length{Kind: blen.NONE, Value: 1}}
	InfoMap[0x22] = Info{"PAN / EF_ID", blen.Length{Kind: blen.LL, Value: 0}}
	InfoMap[0x27] = Info{"Result-Code", blen.Length{Kind: blen.NONE, Value: 1}}
	InfoMap[0x29] = Info{"Terminal ID", blen.Length{Kind: blen.NONE, Value: 4}}
	InfoMap[0x2A] = Info{"VU-number", blen.Length{Kind: blen.NONE, Value: 15}}
	InfoMap[0x3B] = Info{"AID authorisation-attribute", blen.Length{Kind: blen.NONE, Value: 8}}
	InfoMap[0x3C] = Info{"Additional-data", blen.Length{Kind: blen.LLL, Value: 0}}
	InfoMap[0x49] = Info{"Currency", blen.Length{Kind: blen.NONE, Value: 2}}
	InfoMap[0x87] = Info{"Receipt-number", blen.Length{Kind: blen.NONE, Value: 2}}
	InfoMap[0x88] = Info{"Turnover record number", blen.Length{Kind: blen.NONE, Value: 3}}
	InfoMap[0x8A] = Info{"Card-type", blen.Length{Kind: blen.NONE, Value: 1}}
	InfoMap[0x8B] = Info{"Card-name", blen.Length{Kind: blen.LL, Value: 0}}
	InfoMap[0xF1] = Info{"Text1 line 1", blen.Length{Kind: blen.LL}}
	InfoMap[0xF2] = Info{"Text1 line 2", blen.Length{Kind: blen.LL}}
	InfoMap[0xF3] = Info{"Text1 line 3", blen.Length{Kind: blen.LL}}
	InfoMap[0xF4] = Info{"Text1 line 4", blen.Length{Kind: blen.LL}}
}
