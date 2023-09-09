package camplocalapi

import (
	camplocal "github.com/awlsring/camp/generated/camp_local"

	"context"
)

type ApiKeySecurityHandler struct {
}

func (h *ApiKeySecurityHandler) HandleSmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpApiKeyAuth) (context.Context, error) {
	return ctx, nil
}
