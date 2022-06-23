package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceEndpoints struct {
	ClientSet *kubernetes.Clientset
}

func (s *ServiceEndpoints) Register(r chi.Router) {
	r.Get("/service/{namespace}", s.GetServicesByNs)
}

func (s *ServiceEndpoints) GetServicesByNs(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	services, err := s.ClientSet.CoreV1().Services(ns).List(r.Context(),
	metav1.ListOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(services)
}



