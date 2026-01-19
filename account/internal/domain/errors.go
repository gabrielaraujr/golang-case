package errors

import "errors"

// Domain validation errors
var (
	ErrFullNameRequired  = errors.New("full name is required")
	ErrCPFRequired       = errors.New("CPF is required")
	ErrSalaryRequired    = errors.New("salary is required")
	ErrEmailRequired     = errors.New("email is required")
	ErrBirthDateRequired = errors.New("birth date is required")
)

// Domain business logic errors
var (
	ErrOnlyPendingCanStartAnalysis     = errors.New("only pending proposals can start analysis")
	ErrOnlyAnalyzingCanBeApproved      = errors.New("only analyzing proposals can be approved")
	ErrOnlyPendingOrAnalyzingCanReject = errors.New("only pending or analyzing proposals can be rejected")
)

// Domain repository errors
var (
	ErrProposalNotFound = errors.New("proposal not found")
)
