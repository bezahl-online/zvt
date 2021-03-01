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
