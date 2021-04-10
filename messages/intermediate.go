package messages

var IntermediateStatus map[byte]string = make(map[byte]string)

func init() {
	IntermediateStatus[0x00] = "BZT wartet auf Betragbestätigung"
	IntermediateStatus[0x01] = "Bitte Anzeigen auf dem PIN-Pad beachten"
	IntermediateStatus[0x02] = "Bitte Anzeigen auf dem PIN-Pad beachten"
	IntermediateStatus[0x03] = "Vorgang nicht möglich"
	IntermediateStatus[0x04] = "BZT wartet auf Antwort vom FEP"
	IntermediateStatus[0x05] = "BZT sendet Autostorno"
	IntermediateStatus[0x06] = "BZT sendet Nachbuchungen"
	IntermediateStatus[0x07] = "Karte nicht zugelassen"
	IntermediateStatus[0x08] = "Karte unbekannt / undefiniert"
	IntermediateStatus[0x09] = "Karte verfallen"
	IntermediateStatus[0x0A] = "Karte einstecken"
	IntermediateStatus[0x0B] = "Bitte Karte entnehmen!"
	IntermediateStatus[0x0C] = "Karte nicht lesbar"
	IntermediateStatus[0x0D] = "Vorgang abgebrochen"
	IntermediateStatus[0x0E] = "Vorgang wird bearbeitet bitte warten..."

	// IntermediateStatus[0x6C] = "Vorgang wird bearbeitet bitte warten..."
}
