package camp_reporting

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
)

func (c *CampLocalReporting) Heartbeat(ctx context.Context, id machine.Identifier) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("heartbeating system")

	in := &local.HeartbeatInput{
		InternalIdentifier: id.String(),
	}

	log.Debug().Msg("sending heartbeat request")
	_, err := c.client.Heartbeat(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to send heartbeat request")
		return err
	}

	log.Debug().Msg("heartbeat request sent")
	return nil
}
