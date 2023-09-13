package service

import (
	"errors"
	"fmt"

	camplocal "github.com/awlsring/camp/packages/camp_local"

	"context"
)

func SecurityHandler(campKey string, agentKeys []string) camplocal.SecurityHandler {
	return &LocalSecurityHandler{
		campKey:  campKey,
		nodeKeys: agentKeys,
	}
}

type LocalSecurityHandler struct {
	campKey  string
	nodeKeys []string
}

func (h *LocalSecurityHandler) HandleSmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpApiKeyAuth) (context.Context, error) {
	valid := false
	key := t.GetAPIKey()

	switch operationName {
	case "DescribeMachine", "ListMachines", "Health":
		if key == h.campKey {
			valid = true
		}
	case "Heartbeat", "Register", "ReportStatusChange", "ReportSystemChange":
		for _, nodeKey := range h.nodeKeys {
			if key == nodeKey {
				valid = true
			}
		}
	}
	if valid {
		return ctx, nil
	}
	return ctx, errors.New(fmt.Sprintf("Invalid API Key for Operation: %s", operationName))
}
