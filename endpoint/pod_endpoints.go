package endpoint

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodEndpoints struct {
	ClientSet kubernetes.Interface
}

func (p *PodEndpoints) Register(r chi.Router) {
	r.Get("/", p.ListPods)
	r.Get("/pod/{namespace}/{name}", p.GetPod)
}

func (p *PodEndpoints) ListPods(w http.ResponseWriter, r *http.Request) {
	pods, err := p.ClientSet.CoreV1().Pods(apiv1.NamespaceDefault).List(r.Context(),
		metav1.ListOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(pods)
}

func (p *PodEndpoints) GetPod(w http.ResponseWriter, r *http.Request) {
	ns := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")
	pod, err := p.ClientSet.CoreV1().Pods(ns).Get(r.Context(), name, 
	metav1.GetOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(pod)
}
