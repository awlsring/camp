package camplocalapi

import (
	camplocal "github.com/awlsring/camp/generated/camp_local"
	"github.com/awlsring/camp/internal/pkg/server"

	"context"
)

func SecurityHandler(auth server.Auth) camplocal.SecurityHandler {
	return &ApiKeySecurityHandler{
		auth: auth,
	}
}

type ApiKeySecurityHandler struct {
	auth server.Auth
}

func (h *ApiKeySecurityHandler) HandleSmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpApiKeyAuth) (context.Context, error) {
	return h.auth.Authenticate(ctx, operationName, t.GetAPIKey())
}

func (h *ApiKeySecurityHandler) HandleSmithyAPIHttpBearerAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpBearerAuth) (context.Context, error) {
	return h.auth.Authenticate(ctx, operationName, t.GetToken())
}
