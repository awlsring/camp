package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/awlsring/camp/apps/local/machine"
	camplocal "github.com/awlsring/camp/packages/camp_local"
	"github.com/rs/zerolog/log"
)

// Compile-time check for Handler.
var _ camplocal.Handler = (*Handler)(nil)

type Handler struct {
	machine machine.Controller
}

func NewHandler(m machine.Controller) camplocal.Handler {
	return Handler{
		machine: m,
	}
}

func (h Handler) Health(ctx context.Context) (camplocal.HealthRes, error) {
	log.Debug().Msg("Invoke Health")
	return &camplocal.HealthResponseContent{
		Success: true,
	}, nil
}

func (h Handler) DescribeMachine(ctx context.Context, req camplocal.DescribeMachineParams) (camplocal.DescribeMachineRes, error) {
	log.Debug().Msg("Invoke Handler.DescribeMachine")
	m, err := h.machine.DescribeMachine(ctx, req.Identifier)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Debug().Msgf("Machine with identifier %s not found", req.Identifier)
			return &camplocal.ResourceNotFoundExceptionResponseContent{
				Message: fmt.Sprintf("Machine with identifier %s not found", req.Identifier),
			}, nil
		}
		log.Error().Err(err).Msgf("Failed to describe machine with identifier %s", req.Identifier)
		return nil, err
	}

	return &camplocal.DescribeMachineResponseContent{
		Summary: modelToSummary(m),
	}, nil
}

func (h Handler) ListMachines(ctx context.Context) (camplocal.ListMachinesRes, error) {
	log.Debug().Msg("Invoke ListMachines")
	m, err := h.machine.ListMachines(ctx, nil)
	if err != nil {
		return nil, err
	}

	var summaries []camplocal.MachineSummary
	for _, machine := range m {
		summaries = append(summaries, modelToSummary(machine))
	}

	return &camplocal.ListMachinesResponseContent{
		Summaries: summaries,
	}, nil
}

func (h Handler) Heartbeat(ctx context.Context, req *camplocal.HeartbeatRequestContent) (camplocal.HeartbeatRes, error) {
	log.Debug().Msg("Invoke Heartbeat")
	err := h.machine.AcknowledgeHeartbeat(ctx, req.InternalIdentifier)
	if err != nil {
		return nil, err
	}
	return &camplocal.HeartbeatResponseContent{
		Success: true,
	}, nil
}

func (h Handler) Register(ctx context.Context, req *camplocal.RegisterRequestContent) (camplocal.RegisterRes, error) {
	log.Debug().Msg("Invoke Register")
	model := reportedMachineSummaryToModel(&req.Summary)
	err := h.machine.RegisterMachine(ctx, model)
	if err != nil {
		return nil, err
	}
	return &camplocal.RegisterResponseContent{
		Success: true,
	}, nil
}

func (h Handler) ReportStatusChange(ctx context.Context, req *camplocal.ReportStatusChangeRequestContent) (camplocal.ReportStatusChangeRes, error) {
	log.Debug().Msg("Invoke ReportStatusChange")
	err := h.machine.UpdateStatus(ctx, req.InternalIdentifier, machine.MachineStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &camplocal.ReportStatusChangeResponseContent{
		Success: true,
	}, nil
}

func (h Handler) ReportSystemChange(ctx context.Context, req *camplocal.ReportSystemChangeRequestContent) (camplocal.ReportSystemChangeRes, error) {
	log.Debug().Msg("Invoke ReportSystemChange")
	model := reportedMachineSummaryToModel(&req.Summary)
	err := h.machine.UpdateMachine(ctx, model)
	if err != nil {
		return nil, err
	}
	return &camplocal.ReportSystemChangeResponseContent{
		Success: true,
	}, nil
}
