package camplocalapi

import (
	"context"
	"errors"

	camplocal "github.com/awlsring/camp/generated/camp_local"
)

// Compile-time check for Handler.
var _ camplocal.Handler = (*Handler)(nil)

type Handler struct {
}

func NewHandler() camplocal.Handler {
	return Handler{}
}

func (h Handler) Health(ctx context.Context) (camplocal.HealthRes, error) {
	return &camplocal.HealthResponseContent{
		Success: true,
	}, nil
}

func (h Handler) DescribeMachine(ctx context.Context, req camplocal.DescribeMachineParams) (camplocal.DescribeMachineRes, error) {
	return &camplocal.DescribeMachineResponseContent{
		Summary: camplocal.MachineSummary{
			Identifier: req.Identifier,
		},
	}, nil
}

func (h Handler) ListMachines(ctx context.Context) (camplocal.ListMachinesRes, error) {
	return &camplocal.ValidationExceptionResponseContent{
		Message: "not implemented",
		FieldList: []camplocal.ValidationExceptionField{
			{
				Path:    "not implemented",
				Message: "not implemented",
			},
		},
	}, errors.New("not implemented")
}

func (h Handler) Heartbeat(ctx context.Context, req *camplocal.HeartbeatRequestContent) (camplocal.HeartbeatRes, error) {
	panic("not implemented")
}

func (h Handler) Register(ctx context.Context, req *camplocal.RegisterRequestContent) (camplocal.RegisterRes, error) {
	panic("not implemented")
}

func (h Handler) ReportStatusChange(ctx context.Context, req *camplocal.ReportStatusChangeRequestContent) (camplocal.ReportStatusChangeRes, error) {
	panic("not implemented")
}

func (h Handler) ReportSystemChange(ctx context.Context, req *camplocal.ReportSystemChangeRequestContent) (camplocal.ReportSystemChangeRes, error) {
	panic("not implemented")
}
