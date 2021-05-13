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
	"go.uber.org/zap"
)

var fixedPassword [3]byte = [3]byte{0x12, 0x34, 0x56}

type SingleTotals struct {
	ReceiptNrStart int // 2 Byte BCD Belegnummer Start (N4)
	ReceiptNrEnd   int // 2 Byte BCD Belegnummer Ende (N4)
	CountEC        int // 1 Byte binär Anzahl ec-Karte
	TotalEC        int // 6 Byte BCD Umsatzsumme ec-Karte
	CountJCB       int // 1 Byte binär Anzahl JCB
	TotalJCB       int // 6 Byte BCD Umsatzsumme JCB
	CountEurocard  int // 1 Byte binär Anzahl Eurocard
	TotalEurocard  int // 6 Byte BCD Umsatzsumme Eurocard
	CountAmex      int // 1 Byte binär Anzahl Amex
	TotalAmex      int // 6 Byte BCD Umsatzsumme Amex
	CountVisa      int // 1 Byte binär Anzahl VISA
	TotalVisa      int // 6 Byte BCD Umsatzsumme VISA
	CountDiners    int // 1 Byte binär Anzahl Diners
	TotalDiners    int // 6 Byte BCD Umsatzsumme Diners
	CountOther     int // 1 Byte binär Anzahl übrige Karten
	TotalOther     int // 6 Byte BCD Umsatzsumme  übrige Karten
}

func (s *SingleTotals) Unmarshal(data []byte) {
	s.ReceiptNrStart = int(bcd.ToUint64(data[:2]))
	s.ReceiptNrEnd = int(bcd.ToUint64(data[2:4]))
	s.CountEC = int(data[4])
	s.TotalEC = int(bcd.ToUint64(data[5:11]))
	s.CountJCB = int(data[11])
	s.TotalJCB = int(bcd.ToUint64(data[12:18]))
	s.CountEurocard = int(data[18])
	s.TotalEurocard = int(bcd.ToUint64(data[19:25]))
	s.CountAmex = int(data[25])
	s.TotalAmex = int(bcd.ToUint64(data[26:32]))
	s.CountVisa = int(data[32])
	s.TotalVisa = int(bcd.ToUint64(data[33:39]))
	s.CountDiners = int(data[39])
	s.TotalDiners = int(bcd.ToUint64(data[40:46]))
	s.CountOther = int(data[46])
	s.TotalOther = int(bcd.ToUint64(data[47:53]))
}

type EoDResultData struct {
	TraceNr  int
	Date     string
	Time     string
	Total    int64
	Totals   *SingleTotals
	PrintOut *PrintOut
}
type PrintOut struct {
	Type byte
	Text string
}
type EoDResult struct {
	Error  string
	Result string
	Data   *EoDResultData
}

type EndOfDayResponse struct {
	TransactionResponse
	Transaction *EoDResult
}

// EndOfDay implements inst 06 50
// initiates a end of day process
// ECR can instruct the PT to abort execution of the command
func (p *PT) EndOfDay() error {
	Logger.Info("END OF DAY")
	return p.SendCommand(Command{
		CtrlField: instr.Map["EndOfDay"],
		Data: apdu.DataUnit{
			Data: []byte(fixedPassword[:]),
		},
	})
}

func (r *EndOfDayResponse) Process(result *Command) error {
	if r.Transaction == nil {
		r.Transaction = &EoDResult{
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
			case 0x77:
				Logger.Info("not possible")
				r.Transaction.Result = Result_Abort
			}
			return nil
		case 0x0F:
			Logger.Info("Transaktion erfolgreich")
			r.Transaction.Result = Result_Success
			// result could be changed next
			r.Transaction.Data.FromOBJs(result.Data.BMPOBJs)
			return nil
		case 0xD3:
			r.Transaction.Data = &EoDResultData{
				PrintOut: &PrintOut{},
			}
			r.Transaction.Data.FromTLV(result.Data.TLVContainer.Objects)
			Logger.Info("Print Block (06 D3)",
				zap.String("Data", r.Transaction.Data.PrintOut.Text))
		default:
			Logger.Error(fmt.Sprintf("PT command '06 %02X' not handled",
				result.CtrlField.Instr))
		}
	case 0x04:
		switch result.CtrlField.Instr {
		case 0x0F:
			r.Transaction.Data = &EoDResultData{}
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
	case 0x80:
		// got ACK from PT
	default:
		Logger.Error(fmt.Sprintf("PT command '%02X %02X' not handled",
			result.CtrlField.Class, result.CtrlField.Instr))
	}
	return nil
}

func (r *EoDResultData) FromTLV(objs []tlv.DataObject) {
	for _, obj := range objs {
		switch obj.TAG[0] {
		case 0x1F:
			switch obj.TAG[1] {
			case 7: // receipt-type
				r.PrintOut.Type = obj.Data[0]
			default:
				Logger.Error(fmt.Sprintf("TLV TAG '1F %0X' not handled",
					obj.TAG[1]))
			}
		case 0x25:
			r.PrintOut.Text = util.GetPureText(string(obj.Data))
		default:
			Logger.Error(fmt.Sprintf("TLV TAG '% 0X' not handled",
				obj.TAG))
		}
	}

}
func (r *EoDResultData) FromOBJs(objs []bmp.OBJ) (result string, error string) {
	for _, obj := range objs {
		switch obj.ID {
		case 0x04:
			amount := bcd.ToUint64(obj.Data)
			r.Total = int64(amount)
		case 0x0B:
			r.TraceNr = int(bcd.ToUint32(obj.Data))
		case 0x0C:
			r.Time = fmt.Sprintf("%06X", obj.Data)
		case 0x0D:
			r.Date = fmt.Sprintf("%04X", obj.Data)
		case 0x27:
			// Error and Result in AuthResult
			switch obj.Data[0] {
			case 0x6C:
				// result = Result_Abort // just don't
			case 0xE0:
				result = Result_SoftwareUpdate
			}
		case 0x60:
			r.Totals = &SingleTotals{}
			r.Totals.Unmarshal(obj.Data)
		default:
			Logger.Error(fmt.Sprintf("no path for BMP-ID %0X", obj.ID))
		}
	}
	return result, error
}
