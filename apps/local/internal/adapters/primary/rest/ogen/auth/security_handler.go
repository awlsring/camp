package auth

import (
	"context"
	"errors"
	"fmt"

	camplocal "github.com/awlsring/camp/packages/camp_local"
)

type StaticKeySecurityHandler struct {
	apiKeys   []string
	agentKeys []string
}

func NewSecurityHandler(campKey []string, agentKeys []string) camplocal.SecurityHandler {
	return &StaticKeySecurityHandler{ // dynamically load these keys, later
		apiKeys:   campKey,
		agentKeys: agentKeys,
	}
}

func (h *StaticKeySecurityHandler) HandleSmithyAPIHttpApiKeyAuth(ctx context.Context, operationName string, t camplocal.SmithyAPIHttpApiKeyAuth) (context.Context, error) {
	valid := false
	key := t.GetAPIKey()

	switch operationName {
	case "DescribeMachine", "ListMachines", "AddTagsToMachine", "RemoveTagFromMachine", "Health":
		for _, k := range h.apiKeys {
			if key == k {
				valid = true
			}
		}
	case "Heartbeat", "Register", "ReportStatusChange", "ReportSystemChange":
		for _, k := range h.agentKeys {
			if key == k {
				valid = true
			}
		}
	}
	if valid {
		return ctx, nil
	}
	return ctx, errors.New(fmt.Sprintf("Invalid API Key for Operation: %s", operationName))
}
