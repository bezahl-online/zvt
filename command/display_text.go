package command

import (
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
)

// DisplayText implements instr 06 E0
func (p *PT) DisplayText(text []string) error {
	i := instr.Map["DisplayText"]
	return p.send(Command{
		CtrlField: i,
		Data:      compileText(text),
	})
}

func compileText(textarray []string) apdu.DataUnit {
	dataUnit := apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{},
	}
	for i, text := range textarray {
		bmp := bmp.OBJ{
			ID:   0xf0 + byte(i+1),
			Data: []byte{},
		}
		bmp.Data = append(bmp.Data, []byte(text)...) // append bytes of text line
		dataUnit.BMPOBJs = append(dataUnit.BMPOBJs, bmp)
	}
	return dataUnit
}
