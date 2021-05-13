package command

import (
	"fmt"
	"strings"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/messages"
	"github.com/bezahl-online/zvt/util"
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
	Logger.Info(fmt.Sprintf("ECR: AUTHORISATION amount: %5.2f", float64(config.Amount)/100))
	return p.SendCommand(Command{instr.Map["Authorisation"], config.marshal()})
}

func (a *AuthConfig) marshal() apdu.DataUnit {
	return apdu.DataUnit{
		BMPOBJs: []bmp.OBJ{
			{ID: 0x04, Data: bcd.FromUint(uint64(a.Amount), 6)},
		},
	}
}

type CardData struct {
	Name   string
	Type   int
	PAN    string
	Tech   int
	SeqNr  int
	Expiry string
}

// payment-type:
// 40 = offline
// 50  =  card  in  terminal  checked  positively,  but  no  Authorisation  carried out
// 60 = online
// 70 = PIN-payment (also possible for EMV-processing, i.e. credit cards,ecTrack2, ecEMV online/offline).If the TLV-container is active, this information can be specified in tag 2F (see chapter TLV-container).
const (
	PaymentType_offline = 0x40
	PaymentType_positiv = 0x50
	PaymentType_online  = 0x60
	PaymentType_pin     = 0x70
)

type AuthResultData struct {
	Amount      int64
	Currency    int
	ReceiptNr   int
	TurnoverNr  int
	TraceNr     int
	Date        string
	Time        string
	TID         string
	VU          string
	AID         string
	Info        string
	PaymentType byte
	EMVCustomer string
	EMVMerchant string
	Card        CardData
}

const (
	Result_Pending        = "pending"
	Result_Success        = "success"
	Result_Abort          = "abort"
	Result_Timeout        = "timeout"
	Result_Need_EoD       = "need_end_of_day"
	Result_SoftwareUpdate = "software_update"
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
		r.Transaction = &AuthResult{
			Error:  "",
			Result: Result_Pending,
		}
	}
	if r.Transaction.Result == "" {
		r.Transaction.Result = Result_Pending
	}
	switch result.CtrlField.Class {
	case 0x06:
		switch result.CtrlField.Instr {
		case 0x1E: // Abort from PT
			Logger.Info("PT: 'Transaktion abgebrochen'")
			r.Transaction.Result = Result_Abort
			var ok bool
			r.Message, ok = messages.ErrorMessage[result.Data.Data[0]]
			if ok {
				Logger.Info(fmt.Sprintf("PT: '%s'", r.Message))
			} else {
				Logger.Error(fmt.Sprintf("06 1E: unmapped error message code %0X", result.Data.Data[0]))
			}
			return nil
		case 0x0F:
			Logger.Info("PT: 'Transaktion erfolgreich'")
			r.Transaction.Result = Result_Success
			return nil
		default:
			Logger.Error(fmt.Sprintf("PT command '06 %02X' not handled",
				result.CtrlField.Instr))
		}
	case 0x04:
		switch result.CtrlField.Instr {
		case 0x0F:
			r.Transaction.Data = &AuthResultData{}
			r.Transaction.Result = Result_Pending
			// Result can be changed by next call
			r.Transaction.Data.FromOBJs(r, result.Data.BMPOBJs)
			r.Transaction.Data.FromTLV(r, result.Data.TLVContainer.Objects)
		case 0xFF:
			if result.Data.Data != nil && len(result.Data.Data) > 0 {
				r.Status = result.Data.Data[0]
				var ok bool
				r.Message, ok = messages.IntermediateStatus[result.Data.Data[0]]
				if ok {
					Logger.Info(fmt.Sprintf("PT: '%s'", r.Message))
				} else {
					Logger.Error(fmt.Sprintf("04 FF: unmapped intermediate status code %0X", result.Data.Data[0]))
				}
			}
			r.Transaction.Data.FromTLV(r, result.Data.TLVContainer.Objects)
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
func (r *AuthResultData) FromTLV(ar *AuthorisationResponse, objs []tlv.DataObject) {
	for _, obj := range objs {
		switch obj.TAG[0] {
		case 0x24:
			ar.Message = util.GetPureText(string(obj.Data))
			Logger.Info(fmt.Sprintf("PT: '%s'", strings.ReplaceAll(ar.Message, "\n", "; ")))
		case 0x46:
			r.EMVCustomer = string(obj.Data)
		case 0x47:
			r.EMVMerchant = string(obj.Data)
		case 0x1F:
			switch obj.TAG[1] { //FIXME: add all tags
			case 0x10: // FIXME: "4.2.3.Karteninhaberauthentifizierung"
				// r.Card.Auth = int(obj.Data[0])
			case 0x12:
				r.Card.Tech = int(obj.Data[0])
			default:
				Logger.Error(fmt.Sprintf("04 FF TLV TAG %02X' not handled", obj.TAG))
			}
		default:
			Logger.Error(fmt.Sprintf("04 FF TLV TAG %02X' not handled", obj.TAG[0]))
		}
	}
}
func (r *AuthResultData) FromOBJs(ar *AuthorisationResponse, objs []bmp.OBJ) (error string) {
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
		case 0x0E:
			r.Card.Expiry = fmt.Sprintf("%04X", obj.Data)
		case 0x17:
			r.Card.SeqNr = int(bcd.ToUint16(obj.Data))
		case 0x19:
			r.PaymentType = obj.Data[0]
		case 0x22:
			pan := formatPAN(obj.Data)
			r.Card.PAN = pan
		case 0x27:
			// Error and Result in AuthResult
			switch obj.Data[0] {
			case 0x6C:
				// ar.Transaction.Result = Result_Abort // just don't!
			case 0xF0:
				ar.Transaction.Result = Result_Need_EoD
			default:
				ar.Status = obj.Data[0]
				var ok bool
				ar.Message, ok = messages.ErrorMessage[obj.Data[0]]
				if !ok {
					Logger.Error(fmt.Sprintf("0x27: unmapped error message code %0X", obj.Data[0]))
				}
			}
		case 0x29:
			r.TID = fmt.Sprintf("%X", obj.Data)
		case 0x2A:
			r.VU = strings.TrimSpace(string(obj.Data))
		case 0x3B:
			r.AID = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x3C:
			r.Info = string(obj.Data)
			Logger.Info(fmt.Sprintf("PT: Info '%s'", r.Info))
		case 0x49:
			r.Currency = int(bcd.ToUint16(obj.Data))
		case 0x87:
			r.ReceiptNr = int(bcd.ToUint16(obj.Data))
		case 0x88:
			r.TurnoverNr = int(bcd.ToUint64(obj.Data))
		case 0x8B:
			r.Card.Name = strings.Trim(string(obj.Data), string(byte(0x00)))
		case 0x8A:
			r.Card.Type = int(obj.Data[0])
		default:
			Logger.Error(fmt.Sprintf("no path for BMP-ID %0X", obj.ID))
		}
	}
	return error
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
