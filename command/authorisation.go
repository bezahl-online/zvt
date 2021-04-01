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

// AuthConfig is the auth data struct
type AuthConfig struct {
	Amount      int64
	PaymentType byte
	Currency    int
	// TLV         *tlv.Container
}

// Authorisation implents 06 01
// initiates a payment process
// ECR can instruct the PT to abort execution of the command
func (p *PT) Authorisation(config *AuthConfig) error {
	ctrlField := instr.Map["Authorisation"]
	err := p.send(Command{ctrlField, config.marshal()})
	if err != nil {
		return err
	}
	response, err := PaymentTerminal.ReadResponse()
	if err == nil && !response.IsAck() {
		err = fmt.Errorf("error code %0X %0X", response.CtrlField.Class, response.CtrlField.Instr)
	}
	return err
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
		},
	}
}

type CardData struct {
	Name  string
	Type  int
	PAN   string
	Tech  int
	SeqNr int
}

type AuthResultData struct {
	Amount     int64
	ReceiptNr  int
	TurnoverNr int
	TraceNr    int
	Date       string
	Time       string
	TID        string
	VU         string
	AID        string
	Card       CardData
}

const (
	Result_Pending = "pending"
	Result_Success = "success"
	Result_Abort   = "abort"
	Result_Timeout = "timeout"
)

type AuthResult struct {
	Error  string
	Result string
	Data   *AuthResultData
}

type AuthorisationResponse struct {
	TransactionResponse
	Transaction *AuthResult
}

func (r *AuthorisationResponse) Process(result *Command) error {
	if r.Transaction == nil {
		r.Transaction = &AuthResult{}
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
			r.Transaction.Data = &AuthResultData{}
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

func (r *AuthResultData) FromOBJs(objs []bmp.OBJ) (result string, error string) {
	for _, obj := range objs {
		switch obj.ID {
		case 0x04:
			amount := bcd.ToUint64(obj.Data)
			r.Amount = int64(amount)
		case 0x0B:
			r.TraceNr = int(bcd.ToUint32(obj.Data))
		case 0x0C:
			r.Time = fmt.Sprintf("%06X", obj.Data)
		case 0x0D:
			r.Date = fmt.Sprintf("%04X", obj.Data)
		case 0x17:
			r.Card.SeqNr = int(bcd.ToUint16(obj.Data))
		case 0x22:
			pan := formatPAN(obj.Data)
			r.Card.PAN = pan
		case 0x27:
			// Error and Result in AuthResult
			switch obj.Data[0] {
			case 0x6C:
				result = Result_Abort
			}
		case 0x29:
			r.TID = fmt.Sprintf("%X", obj.Data)
		case 0x2A:
			r.VU = strings.TrimSpace(string(obj.Data))
		case 0x3B:
			r.AID = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x87:
			r.ReceiptNr = int(bcd.ToUint16(obj.Data))
		case 0x88:
			r.TurnoverNr = int(bcd.ToUint64(obj.Data))
		case 0x8B:
			r.Card.Name = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x8A:
			r.Card.Type = int(obj.Data[0])
		default:
			fmt.Printf("no path for BMP-ID %0X", obj.ID)
		}
	}
	return result, error
}

func formatPAN(rawPAN []byte) string {
	raw := fmt.Sprintf("%X", rawPAN)
	raw = strings.ReplaceAll(raw, "E", "X")
	pan := raw[0:4]
	l := int(len(raw) / 4)
	for i := 1; i < l; i++ {
		p := i * 4
		pan += " " + raw[p:p+4]
	}
	pan = strings.TrimRight(pan, "F")
	return pan
}
