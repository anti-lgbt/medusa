package daemons

import (
	"github.com/anti-lgbt/medusa/jobs"
)

type CronJobWorker struct {
	Jobs []jobs.Jobs
}

func NewCronJobWorker() *CronJobWorker {
	return &CronJobWorker{
		Jobs: []jobs.Jobs{
			jobs.NewTrendingMusicJob(),
		},
	}
}

func (d *CronJobWorker) Run() {
	for _, job := range d.Jobs {
		d.Process(job)
	}
}

func (d *CronJobWorker) Process(job jobs.Jobs) {
	for {
		job.Process()
	}
}
