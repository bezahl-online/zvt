package messages

var ErrorMessage map[byte]string = make(map[byte]string)

func init() {
	ErrorMessage[0x00] = "OK"
	ErrorMessage[0x01] = "ABGELEHNT"
	ErrorMessage[0x02] = "ABGELEHNT"
	ErrorMessage[0x03] = "FALSCHE VUNR"
	ErrorMessage[0x04] = "ABGELEHNT"
	ErrorMessage[0x05] = "ABGELEHNT"
	ErrorMessage[0x06] = "KARTE UNGUELTIG"
	ErrorMessage[0x07] = "ABGELEHNT"
	ErrorMessage[0x08] = "ABGELEHNT"
	ErrorMessage[0x09] = "FLEET ID FALSCH"
	ErrorMessage[0x0A] = "TRANSAKTION NICHT ERLAUBT"
	ErrorMessage[0x0B] = "ABGELEHNT"
	ErrorMessage[0x0C] = "ABGELEHNT"
	ErrorMessage[0x0D] = "BETRAG UNGUELTIG"
	ErrorMessage[0x0E] = "KARTE UNGUELTIG"
	ErrorMessage[0x0F] = "SYSTEMFEHLER"
	ErrorMessage[0x10] = "ABGELEHNT"
	ErrorMessage[0x11] = "FLEET ID ZU OFT FALSCH"

	ErrorMessage[0x18] = "SYSTEMFEHLER"
	ErrorMessage[0x19] = "STORNO ABGELEHNT"

	ErrorMessage[0x1D] = "FEHLGESCHLAGEN"
	ErrorMessage[0x1E] = "FORMATFEHLER"

	ErrorMessage[0x22] = "ABGELEHNT"

	ErrorMessage[0x29] = "ABGELEHNT"

	ErrorMessage[0x2B] = "ABGELEHNT"

	ErrorMessage[0x32] = "ABGELEHNT"
	ErrorMessage[0x33] = "ABGELEHNT"

	ErrorMessage[0x36] = "KARTE ABGELAUFEN"
	ErrorMessage[0x37] = "PIN FALSCH"
	ErrorMessage[0x38] = "PIN FALSCH"
	ErrorMessage[0x39] = "KARTE FALSCH"
	ErrorMessage[0x3A] = "SYSTEMFEHLER"

	ErrorMessage[0x3F] = "SYSTEMFEHLER"

	ErrorMessage[0x40] = "BETRAG FALSCH"

	ErrorMessage[0x4C] = "SYSTEMFEHLER"

	ErrorMessage[0x50] = "KARTE ABGELAUFEN"
	ErrorMessage[0x51] = "PIN/KPN FALSCH"
	ErrorMessage[0x52] = "KPN FALSCH"
	ErrorMessage[0x53] = "ABGELEHNT"

	ErrorMessage[0x55] = "ABGELEHNT"
	ErrorMessage[0x56] = "TERMINAL INAKTIV"
	ErrorMessage[0x57] = "INITIALISIERUNG"
	ErrorMessage[0x58] = "ABGELEHNT"
	ErrorMessage[0x59] = "TERMINAL GESPERRT"

	ErrorMessage[0x64] = "card not readable"
	ErrorMessage[0x65] = "card-data not present"
	ErrorMessage[0x66] = "processing-error"
	ErrorMessage[0x67] = "function not permitted for ec-and Maestro-cards"
	ErrorMessage[0x68] = "function not permitted for credit-and tank-cards"

	ErrorMessage[0x6A] = "turnover-file full"
	ErrorMessage[0x6B] = "function deactivated (PT not registered)"
	ErrorMessage[0x6C] = "abort via timeout or abort-key"

	ErrorMessage[0x6E] = "card in blocked-list"
	ErrorMessage[0x6F] = "wrong currency"

	ErrorMessage[0x71] = "credit not sufficient (chip-card)"
	ErrorMessage[0x72] = "chip error"
	ErrorMessage[0x73] = "card-data incorrect"
	ErrorMessage[0x74] = "DUKPT engine exhausted"
	ErrorMessage[0x75] = "text not authentic"
	ErrorMessage[0x76] = "PAN not in white list"
	ErrorMessage[0x77] = "end-of-day batch not possible"
	ErrorMessage[0x78] = "card expired"
	ErrorMessage[0x79] = "card not yet valid"
	ErrorMessage[0x7A] = "card unknown"
	ErrorMessage[0x7B] = "fallback to magnetic stripe for girocard not possible"
	ErrorMessage[0x7C] = "fallback to magnetic stripe not possible"
	ErrorMessage[0x7D] = "communication error (no answer)"
	ErrorMessage[0x7E] = "fallback to magnetic stripe not possible, debit advice possible"

	ErrorMessage[0x83] = "function not possible"

	ErrorMessage[0x85] = "key missing"

	ErrorMessage[0x89] = "PIN-pad defective"

	ErrorMessage[0x9A] = "ZVT protocol error"
	ErrorMessage[0x9B] = "error from dial-up/communication fault"
	ErrorMessage[0x9C] = "please wait"

	ErrorMessage[0xA0] = "receiver notready"
	ErrorMessage[0xA1] = "remote station does not respond"
	ErrorMessage[0xA3] = "no connection"
	ErrorMessage[0xA4] = "submission of Geldkarte not possible"
	ErrorMessage[0xA5] = "function not allowed due to PCI-DSS/P2PE rules"
	ErrorMessage[0xA5] = "function not allowed due to PCI-DSS/P2PE rules"

	ErrorMessage[0xB1] = "memory full"
	ErrorMessage[0xB2] = "merchant-journal full"

	ErrorMessage[0xB4] = "already reversed"
	ErrorMessage[0xB5] = "reversal not possible"

	ErrorMessage[0xB7] = "pre-authorisation incorrect (amount too high wrong)"
	ErrorMessage[0xB8] = "error pre-authorisation"

	ErrorMessage[0xBF] = "voltage supply to low"

	ErrorMessage[0xC0] = "card locking mechanism defective"
	ErrorMessage[0xC1] = "merchant-card locked"
	ErrorMessage[0xC2] = "diagnosis required"
	ErrorMessage[0xC3] = "maximum amount exceeded"
	ErrorMessage[0xC4] = "card-profile invalid (new card-profiles must be loaded)"
	ErrorMessage[0xC5] = "payment method not supported"
	ErrorMessage[0xC6] = "currency not applicable"

	ErrorMessage[0xC8] = "amount too small"
	ErrorMessage[0xC9] = "max. transaction-amount toosmall"

	ErrorMessage[0xCB] = "function only allowed in EURO"
	ErrorMessage[0xCC] = "printer not ready"
	ErrorMessage[0xCD] = "Cashback not possible"

	ErrorMessage[0xD2] = "function not permitted for service-cards/bank-customer-cards"
	ErrorMessage[0xDC] = "card inserted"
	ErrorMessage[0xDD] = "error during card-eject "
	ErrorMessage[0xDE] = "error during card-insertion"

	ErrorMessage[0xE0] = "remote-maintenance activated"

	ErrorMessage[0xE2] = "card-reader does not answer / card-reader defective"
	ErrorMessage[0xE3] = "shutter closed"
	ErrorMessage[0xE4] = "Terminal activation required"

	ErrorMessage[0xE7] = "min. one goods-group not found"
	ErrorMessage[0xE8] = "no goods-groups-table loaded"
	ErrorMessage[0xE9] = "restriction-code not permitted"
	ErrorMessage[0xEA] = "card-code not permitted"
	ErrorMessage[0xEB] = "function not executable (PIN-algorithm unknown)"
	ErrorMessage[0xEC] = "PIN-processing not possible"
	ErrorMessage[0xED] = "PIN-pad defective"

	ErrorMessage[0xF0] = "open end-of-day batch present"
	ErrorMessage[0xF1] = "ec-cash/Maestro offline error"
	ErrorMessage[0xF5] = "OPT-error"
	ErrorMessage[0xF6] = "OPT-data not available (= OPT personalisation required)"
	ErrorMessage[0xFA] = "error transmitting offline-transactions (clearing error)"
	ErrorMessage[0xFB] = "turnover data-set defective"
	ErrorMessage[0xFC] = "necessary device not present or defective"
	ErrorMessage[0xFD] = "baudrate not supported"
	ErrorMessage[0xFE] = "register unknown"
	ErrorMessage[0xFF] = "system error (= other/unknown error)"

}
