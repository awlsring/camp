package camp_reporting

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func statusCodeFromDomain(status power.StatusCode) local.StatusCode {
	switch status {
	case power.StatusCodePending:
		return local.StatusCode_PENDING
	case power.StatusCodeRebooting:
		return local.StatusCode_REBOOTING
	case power.StatusCodeRunning:
		return local.StatusCode_RUNNING
	case power.StatusCodeStarting:
		return local.StatusCode_STARTING
	case power.StatusCodeStopping:
		return local.StatusCode_STOPPING
	default:
		return local.StatusCode_STATUSCODE_UNKNOWN
	}
}

func (s *CampLocalReporting) ReportStatus(ctx context.Context, id machine.Identifier, status power.StatusCode) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("reporting status")

	log.Debug().Msg("converting status code")
	statusCode := statusCodeFromDomain(status)

	in := &local.ReportStatusChangeInput{
		InternalIdentifier: id.String(),
		Status:             statusCode,
	}

	log.Debug().Msg("sending report status request")
	_, err := s.client.ReportStatusChange(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to send report status request")
		return err
	}

	log.Debug().Msg("report status request sent")
	return nil
}
