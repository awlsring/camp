package auth

import "context"

type Authorizer interface {
	Authorize(ctx context.Context, operation string, key string, resource string) (context.Context, error)
}

type CasbinAuthorizer struct{}

func (c CasbinAuthorizer) Authorize(ctx context.Context, operation string, key string, resource string) (context.Context, error) {
	return ctx, nil
}
