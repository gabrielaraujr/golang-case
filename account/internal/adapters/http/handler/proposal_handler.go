package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/gabrielaraujr/golang-case/account/internal/application"
	"github.com/gabrielaraujr/golang-case/account/internal/application/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createProposalExecutor interface {
	Execute(ctx context.Context, req *dto.CreateProposalRequest) (*dto.ProposalResponse, error)
}

type getProposalExecutor interface {
	Execute(ctx context.Context, id uuid.UUID) (*dto.ProposalResponse, error)
}

type ProposalHandler struct {
	createUseCase createProposalExecutor
	getUseCase    getProposalExecutor
}

func NewProposalHandler(
	createUseCase createProposalExecutor,
	getUseCase getProposalExecutor,
) *ProposalHandler {
	return &ProposalHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
	}
}

func (h *ProposalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	response, err := h.createUseCase.Execute(r.Context(), &req)
	if err != nil {
		handleApplicationError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *ProposalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "invalid proposal ID")
		return
	}

	response, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		handleApplicationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func handleApplicationError(w http.ResponseWriter, err error) {
	var appErr *appErrors.ApplicationError
	if errors.As(err, &appErr) {
		writeError(w, appErr.StatusCode, appErr.Code, appErr.Error())
		return
	}
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
