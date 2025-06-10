package cronjob

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// Job represents a job to be run.
type Job struct {
	Schedule *CronExpression
	Task     func()
}

// CronScheduler represents a cron job scheduler.
type CronScheduler struct {
	Jobs    []*Job
	mutex   sync.Mutex
	running bool
	stop    chan struct{}
}

// NewCronScheduler creates a new CronScheduler.
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		Jobs: make([]*Job, 0),
	}
}

// AddJob adds a new job to the scheduler.
func (c *CronScheduler) AddJob(expr string, task func()) error {
	schedule, err := ParseCronExpression(expr)
	if err != nil {
		return err
	}
	job := &Job{
		Schedule: schedule,
		Task:     task,
	}
	c.mutex.Lock()
	c.Jobs = append(c.Jobs, job)
	c.mutex.Unlock()
	return nil
}

// RemoveJob removes a job from the scheduler by index.
func (c *CronScheduler) RemoveJob(index int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if index < 0 || index >= len(c.Jobs) {
		return fmt.Errorf("index out of range")
	}
	c.Jobs = append(c.Jobs[:index], c.Jobs[index+1:]...)
	return nil
}

// Start starts the scheduler.
func (c *CronScheduler) Start() {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return
	}
	c.running = true
	if c.stop == nil {
		c.stop = make(chan struct{})
	}
	c.mutex.Unlock()

	go func() {
		for {
			now := time.Now()
			c.mutex.Lock()
			if !c.running {
				c.mutex.Unlock()
				return
			}
			c.mutex.Unlock()

			nextRun := c.timeUntilNextJob(now)
			if nextRun <= 0 {
				// Run due jobs immediately
				c.runDueJobs(now)
				continue
			}
			timer := time.NewTimer(nextRun)
			select {
			case <-timer.C:
				c.runDueJobs(time.Now())
			case <-c.stop:
				timer.Stop()
				return
			}
		}
	}()
}

// Stop stops the scheduler.
func (c *CronScheduler) Stop() {
	c.mutex.Lock()
	if c.running {
		c.running = false
		close(c.stop)
		c.stop = nil
	}
	c.mutex.Unlock()
}

func nextRunTime(expr *CronExpression, fromTime time.Time) time.Time {
	// Start from the next second
	nextTime := fromTime.Add(time.Second - time.Duration(fromTime.Nanosecond()))
	// Limit to prevent infinite loops in case of errors
	maxIterations := 1000000
	for range maxIterations {
		if isTimeMatching(expr, nextTime) {
			return nextTime
		}
		nextTime = nextTime.Add(time.Second)
	}
	// If we exceed maxIterations, return zero time
	return time.Time{}
}

func (c *CronScheduler) runDueJobs(now time.Time) {
	c.mutex.Lock()
	jobsToRun := make([]*Job, 0)
	for _, job := range c.Jobs {
		if isTimeMatching(job.Schedule, now) {
			jobsToRun = append(jobsToRun, job)
		}
	}
	c.mutex.Unlock()

	for _, job := range jobsToRun {
		go func(job *Job) {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic with stack trace
					fmt.Printf("Task panicked: %v\nStack trace:\n%s\n", r, debug.Stack())
				}
			}()
			job.Task()
		}(job)
	}
}

func (c *CronScheduler) timeUntilNextJob(now time.Time) time.Duration {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	minDuration := time.Hour * 24 * 365 // 1 year
	for _, job := range c.Jobs {
		nextRun := nextRunTime(job.Schedule, now)
		if nextRun.IsZero() {
			continue
		}
		duration := nextRun.Sub(now)
		if duration < minDuration {
			minDuration = duration
		}
	}
	return minDuration
}

// ListJobs lists all jobs in the scheduler.
func (c *CronScheduler) ListJobs() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var jobList []string
	for i, job := range c.Jobs {
		jobList = append(jobList, fmt.Sprintf("Job %d: %v", i, job.Schedule))
	}
	return jobList
}

func isTimeMatching(expr *CronExpression, t time.Time) bool {
	if !contains(expr.Seconds, t.Second()) {
		return false
	}
	if !contains(expr.Minutes, t.Minute()) {
		return false
	}
	if !contains(expr.Hours, t.Hour()) {
		return false
	}
	if !contains(expr.DayOfMonth, t.Day()) {
		return false
	}
	if !contains(expr.Month, int(t.Month())) {
		return false
	}
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Adjust for Sunday=0 in Go but 7 in cron
	}
	if !contains(expr.DayOfWeek, weekday%7) {
		return false
	}
	return true
}

func contains(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
