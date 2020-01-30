package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rendau/barot/internal/domain/core"
	"github.com/rendau/barot/internal/interfaces"
)

type API struct {
	lg     interfaces.Logger
	server *http.Server
	cr     *core.St

	lChan chan error
}

func CreateAPI(
	lg interfaces.Logger,
	listen string,
	cr *core.St,
) *API {
	api := &API{
		lg:    lg,
		cr:    cr,
		lChan: make(chan error, 1),
	}

	api.server = &http.Server{
		Addr:         listen,
		Handler:      api.createRouter(),
		ReadTimeout:  20 * time.Minute, //nolint
		WriteTimeout: 20 * time.Second, //nolint
	}

	return api
}

func (a *API) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			a.lChan <- err
		}
	}()
}

func (a *API) Wait() <-chan error {
	return a.lChan
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
