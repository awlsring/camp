package auth

import "context"

type Authenticator interface {
	Authenticate(ctx context.Context, key string) (context.Context, error)
}

type KeyAuthenticator struct{}

func (k KeyAuthenticator) Authenticate(ctx context.Context, key string) (context.Context, error) {
	return ctx, nil
}
