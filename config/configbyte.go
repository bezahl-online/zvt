package config

const (
	// PaymentReceiptPrintedByECR ECR assumes receipt-printout for payment functions
	PaymentReceiptPrintedByECR = 2

	// AdminReceiptPrintedByECR ECR assumes receipt-printout for administration functions
	AdminReceiptPrintedByECR = 4

	// PTSendsIntermediateStatus PTSendsIntermediateStatus
	PTSendsIntermediateStatus = 8

	// AmountInputOnPTNotPossible ECR controls payment function
	AmountInputOnPTNotPossible = 16

	// AdminFunctionOnPTNotPossible ECR controls administration function
	AdminFunctionOnPTNotPossible = 32

	// ECRusingPrintLinesForPrintout ECR print-type
	// unset (0) means: ECR compiles receipts itself from the status-information data
	ECRusingPrintLinesForPrintout = 128
)

const (
	// Service_MenuNOTAssignedToFunctionKey prevents PT from assigning the service menu to the function key
	Service_MenuNOTAssignedToFunctionKey = 1

	// Service_DisplayTextsForCommandsAuthorisationInCAPITALS Pre-initialisation and Reversal will be displayed in capitals
	Service_DisplayTextsForCommandsAuthorisationInCAPITALS = 2
)
