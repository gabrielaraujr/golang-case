package postgres

import (
	"context"
	"errors"

	domainErrors "github.com/gabrielaraujr/golang-case/account/internal/domain"
	"github.com/gabrielaraujr/golang-case/account/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProposalRepository struct {
	db *pgxpool.Pool
}

func NewProposalRepository(db *pgxpool.Pool) *ProposalRepository {
	return &ProposalRepository{db: db}
}

func (r *ProposalRepository) Save(ctx context.Context, proposal *entities.Proposal) error {
	const query = `
		INSERT INTO proposals (
			id,
			full_name,
			cpf,
			email,
			phone,
			birthdate,
			address_street,
			address_city,
			address_state,
			address_zip,
			status,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

	_, err := r.db.Exec(ctx, query,
		proposal.ID,
		proposal.FullName,
		proposal.CPF,
		proposal.Email,
		proposal.Phone,
		proposal.BirthDate,
		proposal.Address.Street,
		proposal.Address.City,
		proposal.Address.State,
		proposal.Address.ZipCode,
		proposal.Status,
		proposal.CreatedAt,
		proposal.UpdatedAt,
	)
	return err
}

func (r *ProposalRepository) Update(ctx context.Context, proposal *entities.Proposal) error {
	const query = `
		UPDATE proposals SET
			status = $2,
			updated_at = $3
		WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query,
		proposal.ID,
		proposal.Status,
		proposal.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return domainErrors.ErrProposalNotFound
	}
	return nil
}

func (r *ProposalRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Proposal, error) {
	const query = `
		SELECT
			id,
			full_name,
			cpf,
			email,
			phone,
			birthdate,
			address_street,
			address_city,
			address_state,
			address_zip,
			status,
			created_at,
			updated_at
		FROM proposals
		WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	return scanProposal(row)
}

func (r *ProposalRepository) FindByCPF(ctx context.Context, cpf string) (*entities.Proposal, error) {
	const query = `
		SELECT
			id,
			full_name,
			cpf,
			email,
			phone,
			birthdate,
			address_street,
			address_city,
			address_state,
			address_zip,
			status,
			created_at,
			updated_at
		FROM proposals
		WHERE cpf = $1`

	row := r.db.QueryRow(ctx, query, cpf)
	return scanProposal(row)
}

func scanProposal(row pgx.Row) (*entities.Proposal, error) {
	var proposal entities.Proposal
	var status string

	err := row.Scan(
		&proposal.ID,
		&proposal.FullName,
		&proposal.CPF,
		&proposal.Email,
		&proposal.Phone,
		&proposal.BirthDate,
		&proposal.Address.Street,
		&proposal.Address.City,
		&proposal.Address.State,
		&proposal.Address.ZipCode,
		&status,
		&proposal.CreatedAt,
		&proposal.UpdatedAt,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, domainErrors.ErrProposalNotFound
	}
	if err != nil {
		return nil, err
	}

	proposal.Status = entities.ProposalStatus(status)
	return &proposal, nil
}
