package http

import (
	"github.com/gabrielaraujr/golang-case/account/internal/adapters/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(proposalHandler *handler.ProposalHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/proposals", func(r chi.Router) {
		r.Post("/", proposalHandler.Create)
		r.Get("/{id}", proposalHandler.GetByID)
	})

	return r
}
