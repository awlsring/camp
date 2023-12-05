package mqtt

import (
	"github.com/awlsring/camp/internal/pkg/worker"
)

type JobDefinition struct {
	Topic          string
	Job            worker.Job
	ConcurrentJobs uint32
}

func NewDefinition(topic string, job worker.Job, concurrentJobs uint32) *JobDefinition {
	return &JobDefinition{
		Topic:          topic,
		Job:            job,
		ConcurrentJobs: concurrentJobs,
	}
}
