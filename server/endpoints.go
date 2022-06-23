package server

import "github.com/go-chi/chi/v5"

type Endpoints interface {
	Register(r chi.Router)
}