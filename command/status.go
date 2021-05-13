package command

import (
	"fmt"

	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/messages"
)

type StatusResultData struct {
	Date   string
	Time   string
	Status byte
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
	return p.SendCommand(Command{CtrlField: instr.Map["Status"],
		Data: apdu.DataUnit{
			Data: d,
		},
	})
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
		case 0x0F:
			Logger.Info("Statusabfrage erfolgreich")
			r.Transaction.Result = Result_Success
			r.Transaction.Data = &StatusResultData{}
			r.Transaction.Data.FromOBJs(result.Data.BMPOBJs)
			r.Message = messages.ErrorMessage[r.Transaction.Data.Status]
			return nil
		default:
			Logger.Error(fmt.Sprintf("PT command '06 %02X' not handled",
				result.CtrlField.Instr))
		}
	case 0x80:
		// got ACK from PT
	default:
		Logger.Error(fmt.Sprintf("PT command '%02X %02X' not handled",
			result.CtrlField.Class, result.CtrlField.Instr))
	}
	return nil
}

func (r *StatusResultData) FromOBJs(objs []bmp.OBJ) (result string, error string) {
	for _, obj := range objs {
		switch obj.ID {
		case 0x0C:
			r.Time = fmt.Sprintf("%06X", obj.Data)
		case 0x0D:
			r.Date = fmt.Sprintf("%04X", obj.Data)
		case 0x27:
			r.Status = obj.Data[0]
		default:
			Logger.Error(fmt.Sprintf("no path for BMP-ID %0X", obj.ID))
		}
	}
	return result, error
}
