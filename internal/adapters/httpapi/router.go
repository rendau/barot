package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (a *API) createRouter() http.Handler {
	r := mux.NewRouter()

	sr := r.PathPrefix("/banners").Subrouter()
	sr.HandleFunc("", a.hBannerAdd).Methods("POST")
	sr.HandleFunc("/{slot_id:[0-9]+}/{banner_id:[0-9]+}", a.hBannerRemove).Methods("DELETE")
	sr.HandleFunc("/select/{slot_id:[0-9]+}/{usr_type_id:[0-9]+}", a.hBannerSelect).Methods("GET")
	sr.HandleFunc("/add_click/{slot_id:[0-9]+}/{banner_id:[0-9]+}/{usr_type_id:[0-9]+}", a.hBannerAddClick).Methods("POST")

	r.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// middleware
	h := http.Handler(r)
	h = cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-Requested-With"},
		MaxAge:         604800,
	}).Handler(h)
	h = a.mwRecovery(h)

	return h
}
