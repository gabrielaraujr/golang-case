package domain

type AnalysisResult struct {
	Approved bool
	Reason   string
}

func NewApproved() AnalysisResult {
	return AnalysisResult{Approved: true, Reason: ""}
}

func NewRejected(reason string) AnalysisResult {
	return AnalysisResult{Approved: false, Reason: reason}
}

func AnalyzeDocuments(payload *ProposalPayload) AnalysisResult {
	if len(payload.CPF) != 11 {
		return NewRejected("CPF must have exactly 11 digits")
	}

	if len(payload.FullName) < 3 {
		return NewRejected("full name must have at least 3 characters")
	}

	return NewApproved()
}

func AnalyzeCredit(payload *ProposalPayload) AnalysisResult {
	const minSalary = 3000.0

	if payload.Salary <= minSalary {
		return NewRejected("salary must be greater than 3000")
	}

	return NewApproved()
}

func AnalyzeFraud(payload *ProposalPayload) AnalysisResult {
	lastDigit := payload.CPF[len(payload.CPF)-1] - '0'

	if lastDigit%2 != 0 {
		return NewRejected("CPF failed fraud check")
	}

	return NewApproved()
}
