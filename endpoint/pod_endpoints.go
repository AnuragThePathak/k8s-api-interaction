package endpoint

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PodEndpoints struct {
}

func (p *PodEndpoints) Register(r chi.Router) {
	r.Get("/", p.ListPods)
}

func (p *PodEndpoints) ListPods(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List pods"))
}
