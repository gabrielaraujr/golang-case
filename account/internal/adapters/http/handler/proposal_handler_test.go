package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestProposalHandler_Create(t *testing.T) {
	t.Run("should return 201 when proposal is created successfully", func(t *testing.T) {
		proposalID := uuid.New()
		expectedResponse := &dto.ProposalResponse{
			ID:        proposalID,
			FullName:  "John Doe",
			CPF:       "12345678901",
			Email:     "john@example.com",
			Phone:     "11999999999",
			BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Address: dto.AddressResponse{
				Street:  "123 Main St",
				City:    "São Paulo",
				State:   "SP",
				ZipCode: "01234-567",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createUseCase := &mockCreateProposalUseCase{
			executeFn: func(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error) {
				return expectedResponse, nil
			},
		}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		reqBody := `{
			"full_name": "John Doe",
			"cpf": "12345678901",
			"email": "john@example.com",
			"phone": "11999999999",
			"birth_date": "15-01-1990",
			"address": {
				"street": "123 Main St",
				"city": "São Paulo",
				"state": "SP",
				"zip_code": "01234-567"
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/proposals", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}

		var response dto.ProposalResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.ID != expectedResponse.ID {
			t.Errorf("expected ID %v, got %v", expectedResponse.ID, response.ID)
		}
		if response.CPF != expectedResponse.CPF {
			t.Errorf("expected CPF %q, got %q", expectedResponse.CPF, response.CPF)
		}
		if response.Status != expectedResponse.Status {
			t.Errorf("expected status %q, got %q", expectedResponse.Status, response.Status)
		}
	})

	t.Run("should return 400 when JSON is invalid", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		req := httptest.NewRequest(http.MethodPost, "/proposals", bytes.NewBufferString("{invalid json"))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var errResponse map[string]string
		json.NewDecoder(rec.Body).Decode(&errResponse)

		if errResponse["code"] != "INVALID_JSON" {
			t.Errorf("expected code INVALID_JSON, got %q", errResponse["code"])
		}
	})

	t.Run("should return 400 when request validation fails", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{
			executeFn: func(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error) {
				return nil, appErrors.NewInvalidInputError(nil)
			},
		}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		reqBody := `{"full_name":"","cpf":"12345678901","email":"j@e.com","phone":"11999999999","birth_date":"15-01-1990","address":{}}`
		req := httptest.NewRequest(http.MethodPost, "/proposals", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("should return 409 when CPF is duplicated", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{
			executeFn: func(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error) {
				return nil, appErrors.NewDuplicateCPFError()
			},
		}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		reqBody := `{"full_name":"John","cpf":"12345678901","email":"j@e.com","phone":"11999999999","birth_date":"15-01-1990","address":{}}`
		req := httptest.NewRequest(http.MethodPost, "/proposals", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", rec.Code)
		}

		var errResponse map[string]string
		json.NewDecoder(rec.Body).Decode(&errResponse)

		if errResponse["code"] != "DUPLICATE_CPF" {
			t.Errorf("expected code DUPLICATE_CPF, got %q", errResponse["code"])
		}
	})

	t.Run("should return 500 when internal error occurs", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{
			executeFn: func(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error) {
				return nil, appErrors.NewInternalError("database error", nil)
			},
		}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		reqBody := `{"full_name":"John","cpf":"12345678901","email":"j@e.com","phone":"11999999999","birth_date":"15-01-1990","address":{}}`
		req := httptest.NewRequest(http.MethodPost, "/proposals", bytes.NewBufferString(reqBody))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})
}

func TestProposalHandler_GetByID(t *testing.T) {
	t.Run("should return 200 when proposal is found", func(t *testing.T) {
		proposalID := uuid.New()
		expectedResponse := &dto.ProposalResponse{
			ID:        proposalID,
			FullName:  "John Doe",
			CPF:       "12345678901",
			Email:     "john@example.com",
			Phone:     "11999999999",
			BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Address: dto.AddressResponse{
				Street:  "123 Main St",
				City:    "São Paulo",
				State:   "SP",
				ZipCode: "01234-567",
			},
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createUseCase := &mockCreateProposalUseCase{}
		getUseCase := &mockGetProposalUseCase{
			executeFn: func(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error) {
				return expectedResponse, nil
			},
		}
		handler := NewProposalHandler(createUseCase, getUseCase)

		req := httptest.NewRequest(http.MethodGet, "/proposals/"+proposalID.String(), nil)
		rec := httptest.NewRecorder()

		routerContext := chi.NewRouteContext()
		routerContext.URLParams.Add("id", proposalID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerContext))

		handler.GetByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var response dto.ProposalResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.ID != expectedResponse.ID {
			t.Errorf("expected ID %v, got %v", expectedResponse.ID, response.ID)
		}
		if response.CPF != expectedResponse.CPF {
			t.Errorf("expected CPF %q, got %q", expectedResponse.CPF, response.CPF)
		}
	})

	t.Run("should return 400 when ID is invalid", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{}
		getUseCase := &mockGetProposalUseCase{}
		handler := NewProposalHandler(createUseCase, getUseCase)

		req := httptest.NewRequest(http.MethodGet, "/proposals/invalid-uuid", nil)
		rec := httptest.NewRecorder()

		routerContext := chi.NewRouteContext()
		routerContext.URLParams.Add("id", "invalid-uuid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerContext))

		handler.GetByID(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var errResponse map[string]string
		json.NewDecoder(rec.Body).Decode(&errResponse)

		if errResponse["code"] != "INVALID_ID" {
			t.Errorf("expected code INVALID_ID, got %q", errResponse["code"])
		}
	})

	t.Run("should return 404 when proposal not found", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{}
		getUseCase := &mockGetProposalUseCase{
			executeFn: func(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error) {
				return nil, appErrors.NewNotFoundError("proposal")
			},
		}
		handler := NewProposalHandler(createUseCase, getUseCase)

		proposalID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/proposals/"+proposalID.String(), nil)
		rec := httptest.NewRecorder()

		routerContext := chi.NewRouteContext()
		routerContext.URLParams.Add("id", proposalID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerContext))

		handler.GetByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}

		var errResponse map[string]string
		json.NewDecoder(rec.Body).Decode(&errResponse)

		if errResponse["code"] != "NOT_FOUND" {
			t.Errorf("expected code NOT_FOUND, got %q", errResponse["code"])
		}
	})

	t.Run("should return 500 when internal error occurs", func(t *testing.T) {
		createUseCase := &mockCreateProposalUseCase{}
		getUseCase := &mockGetProposalUseCase{
			executeFn: func(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error) {
				return nil, appErrors.NewInternalError("database error", nil)
			},
		}
		handler := NewProposalHandler(createUseCase, getUseCase)

		proposalID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/proposals/"+proposalID.String(), nil)
		rec := httptest.NewRecorder()

		routerContext := chi.NewRouteContext()
		routerContext.URLParams.Add("id", proposalID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerContext))

		handler.GetByID(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})
}
