package config

const (
	// PaymentReceiptPrintedByECR ECR assumes receipt-printout for payment functions
	PaymentReceiptPrintedByECR = 2 << iota

	// AdminReceiptPrintedByECR ECR assumes receipt-printout for administration functions
	AdminReceiptPrintedByECR

	// PTSendsIntermediateStatus PTSendsIntermediateStatus
	PTSendsIntermediateStatus

	// AmountInputOnPTpossible ECR controls payment function
	AmountInputOnPTpossible

	// AdminFunctionOnPTpossible ECR controls administration function
	AdminFunctionOnPTpossible
	_

	// ECRusingPrintLinesForPrintout ECR print-type
	ECRusingPrintLinesForPrintout
)

const (
	// ServiceMenuNOTAssignedToFunctionKey prevents PT from assigning the service menu to the function key
	ServiceMenuNOTAssignedToFunctionKey = 1 << iota

	// DisplayTextsForCommandsAuthorisation Pre-initialisation and Reversal will be displayed in capitals
	DisplayTextsForCommandsAuthorisation
)
