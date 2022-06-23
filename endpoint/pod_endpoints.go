package endpoint

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodEndpoints struct {
	ClientSet *kubernetes.Clientset
}

func (p *PodEndpoints) Register(r chi.Router) {
	r.Get("/", p.ListPods)
}

func (p *PodEndpoints) ListPods(w http.ResponseWriter, r *http.Request) {
	pods, err := p.ClientSet.CoreV1().Pods("").List(r.Context(),
		metav1.ListOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d", len(pods.Items))))
}
