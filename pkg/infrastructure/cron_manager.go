package infrastructure

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Schedule string    `json:"schedule"`
	LastRun  time.Time `json:"last_run"`
	NextRun  time.Time `json:"next_run"`
	EntryID  cron.EntryID
}

type CronManager struct {
	cron *cron.Cron
	jobs map[cron.EntryID]*CronJob
	mu   sync.RWMutex
}

func NewCronManager() *CronManager {
	return &CronManager{
		cron: cron.New(cron.WithSeconds()), // Enable seconds field
		jobs: make(map[cron.EntryID]*CronJob),
	}
}

func (c *CronManager) Start() {
	c.cron.Start()
}

func (c *CronManager) Stop() {
	c.cron.Stop()
}

func (c *CronManager) AddJob(name, schedule string, cmd func()) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Wrap cmd to update LastRun
	wrappedCmd := func() {
		// We need to look up the entry to update it
		// This is tricky because the function closes over variables.
		// For simplicity, we won't update LastRun inside the job execution here
		// because we can get Prev/Next from c.cron.Entry(id).
		cmd()
	}

	id, err := c.cron.AddFunc(schedule, wrappedCmd)
	if err != nil {
		return 0, err
	}

	c.jobs[id] = &CronJob{
		ID:       int(id),
		Name:     name,
		Schedule: schedule,
		EntryID:  id,
	}

	return int(id), nil
}

func (c *CronManager) GetJobs() []CronJob {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var list []CronJob
	entries := c.cron.Entries()

	for _, entry := range entries {
		if job, ok := c.jobs[entry.ID]; ok {
			j := *job
			j.LastRun = entry.Prev
			j.NextRun = entry.Next
			list = append(list, j)
		}
	}
	return list
}
func (c *CronManager) GetStatus() map[string]interface{} {
	if c == nil {
		return map[string]interface{}{"active": false, "jobs": []interface{}{}}
	}
	return map[string]interface{}{
		"active": true, // Always true if manager exists
		"jobs":   c.GetJobs(),
	}
}
