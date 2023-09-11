package camplocalapi

import (
	camplocal "github.com/awlsring/camp/generated/camp_local"
	"github.com/awlsring/camp/internal/pkg/server/auth"

	"context"
)

func SecurityHandler(auth auth.Handler) camplocal.SecurityHandler {
	return &LocalSecurityHandler{
		auth: auth,
	}
}

type LocalSecurityHandler struct {
	auth auth.Handler
}

func (h *LocalSecurityHandler) HandleSmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpApiKeyAuth) (context.Context, error) {
	return h.auth.Authenticate(ctx, operationName, t.GetAPIKey())
}

func (h *LocalSecurityHandler) HandleSmithyAPIHttpBearerAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpBearerAuth) (context.Context, error) {
	return h.auth.Authenticate(ctx, operationName, t.GetToken())
}
