package ogen_server

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Rest struct {
	server *http.Server
}

func NewRest(srv *http.Server) *Rest {
	return &Rest{
		server: srv,
	}
}

func (r *Rest) Start(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	wg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return r.Stop(ctx)
	})

	return wg.Wait()
}

func (r *Rest) Stop(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
