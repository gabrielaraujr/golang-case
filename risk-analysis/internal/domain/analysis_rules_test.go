package domain

import "testing"

func TestAnalyzeDocuments(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		fullName string
		want     bool
	}{
		{name: "valid documents", cpf: "12345678902", fullName: "John Doe", want: true},
		{name: "minimal valid name", cpf: "12345678901", fullName: "Joe", want: true},
		{name: "invalid CPF too short", cpf: "123456789", fullName: "John Doe", want: false},
		{name: "invalid CPF too long", cpf: "123456789012", fullName: "John Doe", want: false},
		{name: "invalid name too short", cpf: "12345678902", fullName: "Jo", want: false},
		{name: "empty name", cpf: "12345678902", fullName: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &ProposalPayload{
				CPF:      tt.cpf,
				FullName: tt.fullName,
				Salary:   5000.0,
			}
			result := AnalyzeDocuments(payload)
			if result.Approved != tt.want {
				t.Errorf("AnalyzeDocuments(%q, %q) = %v, want %v", tt.cpf, tt.fullName, result.Approved, tt.want)
			}
		})
	}
}

func TestAnalyzeCredit(t *testing.T) {
	tests := []struct {
		name   string
		salary float64
		want   bool
	}{
		{name: "above threshold", salary: 5000.0, want: true},
		{name: "just above threshold", salary: 3000.01, want: true},
		{name: "at threshold", salary: 3000.0, want: false},
		{name: "below threshold", salary: 2999.99, want: false},
		{name: "zero", salary: 0.0, want: false},
		{name: "negative", salary: -100.0, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &ProposalPayload{
				CPF:      "12345678902",
				FullName: "John Doe",
				Salary:   tt.salary,
			}
			result := AnalyzeCredit(payload)
			if result.Approved != tt.want {
				t.Errorf("AnalyzeCredit(salary=%.2f) = %v, want %v", tt.salary, result.Approved, tt.want)
			}
		})
	}
}

func TestAnalyzeFraud(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
		want bool
	}{
		{name: "even digit 0", cpf: "12345678900", want: true},
		{name: "even digit 2", cpf: "12345678902", want: true},
		{name: "even digit 4", cpf: "12345678904", want: true},
		{name: "even digit 6", cpf: "12345678906", want: true},
		{name: "even digit 8", cpf: "12345678908", want: true},
		{name: "odd digit 1", cpf: "12345678901", want: false},
		{name: "odd digit 3", cpf: "12345678903", want: false},
		{name: "odd digit 5", cpf: "12345678905", want: false},
		{name: "odd digit 7", cpf: "12345678907", want: false},
		{name: "odd digit 9", cpf: "12345678909", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &ProposalPayload{
				CPF:      tt.cpf,
				FullName: "John Doe",
				Salary:   5000.0,
			}
			result := AnalyzeFraud(payload)
			if result.Approved != tt.want {
				t.Errorf("AnalyzeFraud(cpf=%q) = %v, want %v", tt.cpf, result.Approved, tt.want)
			}
		})
	}
}
