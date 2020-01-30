package http_api

import (
	"net/http"
)

func (a *Api) mwRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				a.lg.Errorw(
					"Panic in http handler",
					err,
					"method", r.Method,
					"path", r.URL,
				)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
