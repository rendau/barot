package http_api

import (
	"context"
	"github.com/rendau/barot/internal/domain/core"
	"github.com/rendau/barot/internal/interfaces"
	"net/http"
	"time"
)

type Api struct {
	lg     interfaces.Logger
	server *http.Server
	cr     *core.St

	lChan chan error
}

func CreateApi(
	lg interfaces.Logger,
	listen string,
	cr *core.St,
) *Api {
	api := &Api{
		lg:    lg,
		cr:    cr,
		lChan: make(chan error, 1),
	}

	api.server = &http.Server{
		Addr:         listen,
		Handler:      api.createRouter(),
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Second,
	}

	return api
}

func (a *Api) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			a.lChan <- err
		}
	}()
}

func (a *Api) Wait() <-chan error {
	return a.lChan
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
