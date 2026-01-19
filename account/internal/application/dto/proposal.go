package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddressRequest struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

type CreateProposalRequest struct {
	FullName  string         `json:"full_name"`
	CPF       string         `json:"cpf"`
	Salary    float64        `json:"salary"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	BirthDate string         `json:"birthdate"`
	Address   AddressRequest `json:"address"`
}

type ProposalResponse struct {
	ID        uuid.UUID       `json:"id"`
	FullName  string          `json:"full_name"`
	CPF       string          `json:"cpf"`
	Salary    float64         `json:"salary"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	BirthDate time.Time       `json:"birthdate"`
	Address   AddressResponse `json:"address"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type AddressResponse struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}
