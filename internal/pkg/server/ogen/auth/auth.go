package auth

import "context"

type Handler interface {
	Authorize(ctx context.Context, operationName string, k string) (context.Context, error)
	Authenticate(ctx context.Context, operationName string, k string) (context.Context, error)
}

type AuthHandler struct {
	n Authenticator
	z Authorizer
}

func NewAuthHandler(n Authenticator, z Authorizer) Handler {
	return &AuthHandler{
		n: n,
		z: z,
	}
}

func (a *AuthHandler) Authorize(ctx context.Context, operationName string, k string) (context.Context, error) {
	// figure out how to get request. maybe make auth z in middleware
	return a.z.Authorize(ctx, operationName, k, "camplocal")
}

func (a *AuthHandler) Authenticate(ctx context.Context, operationName string, k string) (context.Context, error) {
	return a.n.Authenticate(ctx, k)
}
