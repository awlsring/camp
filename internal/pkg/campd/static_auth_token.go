package campd

import (
	"context"

	campagent "github.com/awlsring/camp/packages/camp_agent"
)

type StaticAuthKeyProvider struct {
	key campagent.SmithyAPIHttpApiKeyAuth
}

func NewStaticAuthKeyProvider(token string) *StaticAuthKeyProvider {
	key := campagent.SmithyAPIHttpApiKeyAuth{
		APIKey: token,
	}
	return &StaticAuthKeyProvider{
		key: key,
	}
}

func (p *StaticAuthKeyProvider) SmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string) (campagent.SmithyAPIHttpApiKeyAuth, error) {
	return p.key, nil
}
