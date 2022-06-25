package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	openelb "github.com/openelb/openelb/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenelbEndpoints struct {
	Client client.Client
}

func (o *OpenelbEndpoints) Register(r chi.Router) {
	r.Get("/openelb/bgp-conf", o.GetBgpConf)
	r.Get("/openelb/bgp-peers", o.ListBgpPeers)
}

func (o *OpenelbEndpoints) GetBgpConf(w http.ResponseWriter, r *http.Request) {
	var conf openelb.BgpConf
	err := o.Client.Get(r.Context(), client.ObjectKey{Name: "default"}, &conf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(conf)
}

func (o *OpenelbEndpoints) ListBgpPeers(w http.ResponseWriter, r *http.Request) {
	var peers openelb.BgpPeerList
	err := o.Client.List(r.Context(), &peers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(peers)
}
