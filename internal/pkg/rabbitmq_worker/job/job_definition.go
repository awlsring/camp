package job

import (
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/exchange"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/queue"
)

type Definition struct {
	Exchange       *exchange.Definition
	Queue          *queue.Definition
	Job            Job
	ConcurrentJobs uint32
}

func NewDefinition(queue *queue.Definition, exchange *exchange.Definition, job Job, concurrentJobs uint32) *Definition {
	return &Definition{
		Exchange:       exchange,
		Queue:          queue,
		Job:            job,
		ConcurrentJobs: concurrentJobs,
	}
}
