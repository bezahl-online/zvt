package payment

const (
	// UseDataFromPreviousReadCard The PT should execute the payment using the data from the previous „Read Card“ command.
	//If no card-data is available, the PT sets the corresponding return-code in the Status-Information.
	UseDataFromPreviousReadCard = 2

	// PrinterReady (mainly used for evaluation tests)
	PrinterReady = 4

	// TippableTransaction (since DCPOS 2.5: ignored for EMV tip/tippable transactions)
	TippableTransaction = 8

	// Geldkarte for GiroCard (ignored for DC POS realted or other cards)
	Geldkarte = 16

	// OnlineWithoutPIN (OLV or EuroELV, if only EuroELV is supported by PT)
	// (ignored by DC POS related or other cards)
	OnlineWithoutPIN = 32

	// GirocardTransaction according to TA7.0 rules for TA 7.0 capable PTs
	// DC POS transaction for capable PT's otherwise ignored or refused
	// PIN based transaction for other cards
	GirocardTransaction = 48

	// PaymentExcludeGeldKarte Payment according to PTs decision excluding GeldKarte
	PaymentExcludeGeldKarte = 64

	// PaymentIncludeGeldKarte Payment according to PTs decision including GeldKarte
	PaymentIncludeGeldKarte = 65
)
