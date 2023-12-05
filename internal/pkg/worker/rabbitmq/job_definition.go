package rabbitmq

import (
	"github.com/awlsring/camp/internal/pkg/worker"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/exchange"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/queue"
)

type JobDefinition struct {
	Exchange       *exchange.Definition
	Queue          *queue.Definition
	Job            worker.Job
	ConcurrentJobs uint32
}

func NewJobDefinition(queue *queue.Definition, exchange *exchange.Definition, job worker.Job, concurrentJobs uint32) *JobDefinition {
	return &JobDefinition{
		Exchange:       exchange,
		Queue:          queue,
		Job:            job,
		ConcurrentJobs: concurrentJobs,
	}
}
