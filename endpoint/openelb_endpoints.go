package endpoint

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	openelb "github.com/openelb/openelb/api/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenelbEndpoints struct {
	Client client.Client
}

func (o *OpenelbEndpoints) Register(r chi.Router) {
	r.Get("/openelb/bgp-conf", o.GetBgpConf)
	r.Delete("/openelb/bgp-peer", o.DeletePeer)
	r.Get("/openelb/bgp-peer", o.ListBgpPeers)
	r.Get("/openelb/bgp-peer/{name}", o.GetBgpConf)
	r.Post("/openelb/bgp-peer", o.CreateBgpPeer)
	r.Put("/openelb/bgp-peer", o.UpdateBgpPeer)
}

func (o *OpenelbEndpoints) GetBgpConf(w http.ResponseWriter, r *http.Request) {
	var conf openelb.BgpConf
	err := o.Client.Get(r.Context(), client.ObjectKey{Name: "default"}, &conf)
	if err != nil {
		if errors.IsGone(err) {
			log.Println("BgpConf not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
	json.NewEncoder(w).Encode(conf)
}

func (o *OpenelbEndpoints) CreateBgpPeer(w http.ResponseWriter, r *http.Request) {
	var peer openelb.BgpPeer
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
	if err := o.Client.Create(r.Context(), &peer); err != nil {
		if errors.IsAlreadyExists(err) {
			log.Println("BgpPeer already exists")
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
}

func (o *OpenelbEndpoints) DeletePeer(w http.ResponseWriter, r *http.Request) {
	var peer openelb.BgpPeer
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
	if err := o.Client.Delete(r.Context(), &peer); err != nil {
		if errors.IsNotFound(err) {
			log.Println("BgpPeer not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
}

func (o *OpenelbEndpoints) GetBgpPeer(w http.ResponseWriter, r *http.Request) {
	var peer openelb.BgpPeer
	err := o.Client.Get(r.Context(), client.ObjectKey{Name: chi.URLParam(r, "name")}, &peer)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("BgpPeer not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(peer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
}

func (o *OpenelbEndpoints) ListBgpPeers(w http.ResponseWriter, r *http.Request) {
	var peers openelb.BgpPeerList
	err := o.Client.List(r.Context(), &peers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
	err = json.NewEncoder(w).Encode(peers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%T\n", err)
		return
	}
}

func (o *OpenelbEndpoints) UpdateBgpPeer(w http.ResponseWriter, r *http.Request) {
	var peer openelb.BgpPeer
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		log.Println(err)
		return
	}

	if err := o.Client.Update(r.Context(), &peer); err != nil {
		if errors.IsNotFound(err) {
			log.Println("BgpPeer not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.IsConflict(err) {
			log.Println("BgpPeer already exists")
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.IsInvalid(err) {
			log.Println("BgpPeer invalid")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.IsBadRequest(err) {
			log.Println("BgpPeer bad req")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.IsNotAcceptable(err) {
			log.Println("BgpPeer not acceptable")
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Println(errors.ReasonForError(err))
		log.Printf("%T\n", err)
		return
	}
}
