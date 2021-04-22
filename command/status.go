package command

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/messages"
	"github.com/bezahl-online/zvt/util"
)

type StatusResultData struct {
	Date string
	Time string
}

type StatusResult struct {
	Error  string
	Result string
	Data   *StatusResultData
}

type StatusResponse struct {
	TransactionResponse
	Transaction *StatusResult
}

// Status implements inst 05 01
// ECR can request the Status of the PT
func (p *PT) Status() error {
	Logger.Info("STATUS")
	d := []byte(fixedPassword[:])
	d = append(d, 0x03, 0x07) // send TLV
	if err := p.send(Command{CtrlField: instr.Map["Status"],
		Data: apdu.DataUnit{
			Data: d,
		},
	}); err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err != nil {
		return err
	}
	return response.IsAck()
}

func (r *StatusResponse) Process(result *Command) error {
	if r.Transaction == nil {
		r.Transaction = &StatusResult{
			Error:  "",
			Result: Result_Pending,
		}
	}
	switch result.CtrlField.Class {
	case 0x06:
		switch result.CtrlField.Instr {
		case 0x1E:
			r.Status = result.Data.Data[0]
			switch result.Data.Data[0] {
			case 0x6C:
				Logger.Info("Transaktion abgebrochen")
				r.Transaction.Result = Result_Abort
			default:
				r.Message = messages.ErrorMessage[r.Status]
				Logger.Error(fmt.Sprintf("0x1E: no path for result code %0X", result.Data.Data[0]))
			}
			return nil
		case 0x0F:
			Logger.Info("Transaktion erfolgreich")
			r.Transaction.Result = Result_Success
			if result.Data.Data != nil && len(result.Data.Data) > 0 {
				r.Status = result.Data.Data[0]
			}
			return nil
		default:
			Logger.Error(fmt.Sprintf("PT command '06 %02X' not handled",
				result.CtrlField.Instr))
		}
	case 0x04:
		switch result.CtrlField.Instr {
		case 0x0F:
			r.Transaction.Data = &StatusResultData{}
			r.Transaction.Data.FromOBJs(result.Data.BMPOBJs)
			r.Transaction.Result = Result_Pending
			return nil
		case 0xFF:
			if result.Data.Data != nil && len(result.Data.Data) > 0 {
				r.Status = result.Data.Data[0]
			}
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
func (r *StatusResultData) FromOBJs(objs []bmp.OBJ) (result string, error string) {
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
