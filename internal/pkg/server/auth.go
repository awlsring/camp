package server

import "context"

type Auth interface {
	Authorize(ctx context.Context, operationName string, k string) (context.Context, error)
	Authenticate(ctx context.Context, operationName string, k string) (context.Context, error)
}

type ApiKeyAuth struct{}

func NewApiKeyAuth() ApiKeyAuth {
	return ApiKeyAuth{}
}

func (a ApiKeyAuth) Authorize(ctx context.Context, operationName string, k string) (context.Context, error) {
	return ctx, nil
}

func (a ApiKeyAuth) Authenticate(ctx context.Context, operationName string, k string) (context.Context, error) {
	return ctx, nil
}
