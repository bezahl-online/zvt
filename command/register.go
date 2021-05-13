package command

import (
	"fmt"

	"github.com/albenik/bcd"
	"github.com/bezahl-online/zvt/apdu"
	"github.com/bezahl-online/zvt/apdu/bmp"
	"github.com/bezahl-online/zvt/apdu/tlv"
	"github.com/bezahl-online/zvt/config"
	"github.com/bezahl-online/zvt/instr"
	"github.com/bezahl-online/zvt/messages"
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
func (p *PT) Register() error {
	config := configure()
	p.Logger.Info("Register (06 00)")
	return p.SendCommand(Command{
		CtrlField: instr.Map["Registration"],
		Data:      config.CompileConfig(),
	})
}

func configure() Config {
	configByte := config.PaymentReceiptPrintedByECR +
		config.AdminReceiptPrintedByECR +
		config.PTSendsIntermediateStatus +
		config.AmountInputOnPTNotPossible +
		config.AdminFunctionOnPTNotPossible
	serviceByte := config.Service_MenuNOTAssignedToFunctionKey +
		config.Service_DisplayTextsForCommandsAuthorisationInCAPITALS
	// var msgSquID tlv.DataObject = tlv.DataObject{
	// 	TAG:  []byte{0x1F, 0x73},
	// 	Data: []byte{0, 0, 0},
	// }
	// var cardType tlv.DataObject = tlv.DataObject{
	// 	TAG:  []byte{0x1F, 0x60},
	// 	Data: []byte{0x03},
	// }
	var listOfCommands tlv.DataObject = tlv.DataObject{
		TAG:  []byte{0x26},
		Data: []byte{0x0A, 0x02, 0x06, 0xD3},
	}
	var tlvContainer *tlv.Container = &tlv.Container{
		Objects: []tlv.DataObject{},
	}
	tlvContainer.Objects = append(tlvContainer.Objects, listOfCommands) // , msgSquID , cardType)
	return Config{
		pwd:          fixedPassword,
		config:       byte(configByte),
		currency:     EUR,
		service:      byte(serviceByte),
		tlvContainer: tlvContainer,
	}
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
	case 0x80:
		// got ACK from PT
	default:
		Logger.Error(fmt.Sprintf("PT command '%02X %02X' not handled",
			result.CtrlField.Class, result.CtrlField.Instr))
	}
	return nil
}

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
				// result = Result_Abort // just don't
			case 0xF0:
				result = Result_Need_EoD
			default:
				Logger.Error(fmt.Sprintf("0x6C: no path for status %0X", obj.Data[0]))
			}
		default:
			Logger.Error(fmt.Sprintf("no path for BMP-ID %0X", obj.ID))
		}
	}
	return result, error
}
