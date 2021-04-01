package command

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/instr"
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
	TraceNr int
	Date    string
	Time    string
	Total   int64
	Totals  SingleTotals
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
	if err := p.send(Command{
		CtrlField: instr.Map["EndOfDay"],
		Data: apdu.DataUnit{
			Data: []byte(fixedPassword[:]),
		},
	}); err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err == nil && !response.IsAck() {
		err = fmt.Errorf("error code %0X %0X", response.CtrlField.Class, response.CtrlField.Instr)
	}
	return err
}

func (r *EndOfDayResponse) Process(result *Command) error {
	if r.Transaction == nil {
		r.Transaction = &EoDResult{}
	}
	switch result.CtrlField.Class {
	case 0x06:
		switch result.CtrlField.Instr {
		case 0x1E:
			switch result.Data.Data[0] {
			case 0x6C:
				fmt.Println("Transaction aborted")
				r.Transaction.Result = Result_Abort
			}
			return nil

		case 0x0F:
			fmt.Println("Transaction successfull")
			r.Transaction.Result = Result_Success
			return nil
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
					r.Message = strings.Map(func(r rune) rune {
						if (unicode.IsLetter(r) ||
							unicode.IsDigit(r) ||
							unicode.IsPunct(r) ||
							unicode.IsSpace(r)) &&
							r != 0x26 {
							return r
						}
						return -1
					}, string(obj.Data))
				}
			}
		}
	}
	return nil
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
				result = Result_Abort
			}
		case 0x60:
			r.Totals.Unmarshal(obj.Data)
		default:
			fmt.Printf("no path for BMP-ID %0X", obj.ID)
		}
	}
	return result, error
}
