package command

import (
	"fmt"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/messages"
	"github.com/bezahl-online/zvt/util"
)

// Config is the config struct
type Config struct {
	pwd          [3]byte
	config       byte
	currency     int // default EUR
	service      byte
	tlvContainer *tlv.Container
}

type RegisterResultData struct {
	Date string
	Time string
}

type RegisterResult struct {
	Error  string
	Result string
	Data   *RegisterResultData
}

type RegisterResponse struct {
	TransactionResponse
	Transaction *RegisterResult
}

// Register implements inst 06 00
// set up different configurations on the PT
func (p *PT) Register(config *Config) error {
	p.Logger.Info(fmt.Sprintf("Register (06 00) Config %04b %04b", (config.config&0xF0)>>4, config.config&0x0F))
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

func (r *RegisterResponse) Process(result *Command) error {
	if r.Transaction == nil {
		r.Transaction = &RegisterResult{
			Error:  "",
			Result: Result_Pending,
		}
	}
	switch result.CtrlField.Class {
	case 0x06:
		switch result.CtrlField.Instr {
		case 0x1E:
			switch result.Data.Data[0] {
			case 0x6C:
				Logger.Info("Transaktion abgebrochen")
				r.Transaction.Result = Result_Abort
			default:
				Logger.Error(fmt.Sprintf("0x1E: no path for result code %0X", result.Data.Data[0]))
			}
			return nil
		case 0x0F:
			Logger.Info("Transaktion erfolgreich")
			r.Transaction.Result = Result_Success
			return nil
		default:
			Logger.Error(fmt.Sprintf("PT command '06 %02X' not handled",
				result.CtrlField.Instr))
		}
	case 0x04:
		switch result.CtrlField.Instr {
		case 0x0F:
			r.Transaction.Data = &RegisterResultData{}
			r.Transaction.Data.FromOBJs(result.Data.BMPOBJs)
			r.Transaction.Result = Result_Pending
			return nil
		case 0xFF:
			r.Status = result.Data.Data[0]
			for _, obj := range result.Data.TLVContainer.Objects {
				if obj.TAG[0] == byte(0x24) {
					r.Message = util.GetPureText(string(obj.Data))
				}
			}
			if len(r.Message) == 0 {
				r.Message = messages.IntermediateStatus[r.Status]
			}
		default:
			Logger.Error(fmt.Sprintf("PT command '04 %02X' not handled",
				result.CtrlField.Instr))
		}
	default:
		Logger.Error(fmt.Sprintf("PT command '%02X %02X' not handled",
			result.CtrlField.Class, result.CtrlField.Instr))
	}
	return nil
}

// func (r *EoDResultData) FromTLV(objs []tlv.DataObject) {
// 	for _, obj := range objs {
// 		switch obj.TAG[0] {
// 		case 0x1F:
// 			switch obj.TAG[1] {
// 			case 7: // receipt-type
// 				r.PrintOut.Type = obj.Data[0]
// 			default:
// 				Logger.Error(fmt.Sprintf("TLV TAG '1F %0X' not handled",
// 					obj.TAG[1]))
// 			}
// 		case 0x25:
// 			r.PrintOut.Text = util.GetPureText(string(obj.Data))
// 		default:
// 			Logger.Error(fmt.Sprintf("TLV TAG '% 0X' not handled",
// 				obj.TAG))
// 		}
// 	}

// }
func (r *RegisterResultData) FromOBJs(objs []bmp.OBJ) (result string, error string) {
	for _, obj := range objs {
		switch obj.ID {
		case 0x0C:
			r.Time = fmt.Sprintf("%06X", obj.Data)
		case 0x0D:
			r.Date = fmt.Sprintf("%04X", obj.Data)
		case 0x27:
			// Error and Result in AuthResult
			switch obj.Data[0] {
			case 0x6C:
				result = Result_Abort
			default:
				Logger.Error(fmt.Sprintf("0x6C: no path for status %0X", obj.Data[0]))
			}
		default:
			Logger.Error(fmt.Sprintf("no path for BMP-ID %0X", obj.ID))
		}
	}
	return result, error
}
