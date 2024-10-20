package cronjob

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestParseCronExpression tests the parsing of valid and invalid cron expressions.
func TestParseCronExpression(t *testing.T) {
	tests := []struct {
		expr       string
		shouldPass bool
	}{
		{"* * * * *", true},
		{"*/15 * * * *", true},
		{"0 12 * * Mon-Fri", true},
		{"0 0 1 1 *", true},
		{"0 0 1 Jan *", true},
		{"0 0 * * Sun", true},
		{"0 0 * * 0", true},          // Sunday as 0
		{"0 0 * * 6", true},          // Sunday as 6 (sometimes accepted)
		{"60 * * * *", false},        // Invalid minute
		{"* 24 * * *", false},        // Invalid hour
		{"* * 32 * *", false},        // Invalid day of month
		{"* * * 13 *", false},        // Invalid month
		{"* * * * SunFunday", false}, // Invalid day of week
		{"* * *", false},             // Too few fields
	}

	for _, test := range tests {
		_, err := ParseCronExpression(test.expr)
		if test.shouldPass && err != nil {
			t.Errorf("Expected expression '%s' to pass, but got error: %v", test.expr, err)
		}
		if !test.shouldPass && err == nil {
			t.Errorf("Expected expression '%s' to fail, but it passed", test.expr)
		}
	}
}

// TestSchedulerAddJob tests adding jobs to the scheduler.
func TestSchedulerAddJob(t *testing.T) {
	scheduler := NewCronScheduler()

	err := scheduler.AddJob("* * * * *", func() {})
	if err != nil {
		t.Errorf("Failed to add valid job: %v", err)
	}

	err = scheduler.AddJob("invalid cron", func() {})
	if err == nil {
		t.Errorf("Expected error when adding job with invalid cron expression")
	}
}

// TestSchedulerRemoveJob tests removing jobs from the scheduler.
func TestSchedulerRemoveJob(t *testing.T) {
	scheduler := NewCronScheduler()

	_ = scheduler.AddJob("* * * * *", func() {})
	_ = scheduler.AddJob("*/5 * * * *", func() {})

	err := scheduler.RemoveJob(1)
	if err != nil {
		t.Errorf("Failed to remove job: %v", err)
	}

	err = scheduler.RemoveJob(5) // Invalid index
	if err == nil {
		t.Errorf("Expected error when removing job with invalid index")
	}
}

// TestSchedulerExecution tests if the scheduler executes tasks at the correct time.
func TestSchedulerExecution(t *testing.T) {
	scheduler := NewCronScheduler()

	var wg sync.WaitGroup
	wg.Add(1)

	// Schedule a job to run one minute from now
	now := time.Now()
	minute := (now.Minute() + 1) % 60
	cronExpr := fmt.Sprintf("%d %d %d %d %d", minute, now.Hour(), now.Day(), int(now.Month()), now.Weekday())

	err := scheduler.AddJob(cronExpr, func() {
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Failed to add job: %v", err)
	}

	scheduler.Start()
	defer scheduler.Stop()

	// Wait for the job to be executed or timeout after 70 seconds
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Job executed successfully
	case <-time.After(70 * time.Second):
		t.Errorf("Scheduled job did not execute in expected time")
	}
}

// TestIsTimeMatching tests the time matching function directly.
func TestIsTimeMatching(t *testing.T) {
	exprStr := "15 14 1 1 *" // At 14:15 on the 1st of January
	expr, err := ParseCronExpression(exprStr)
	if err != nil {
		t.Fatalf("Failed to parse cron expression: %v", err)
	}

	matchingTime := time.Date(2023, time.January, 1, 14, 15, 0, 0, time.UTC)
	if !isTimeMatching(expr, matchingTime) {
		t.Errorf("Expected time %v to match expression %s", matchingTime, exprStr)
	}

	nonMatchingTime := time.Date(2023, time.January, 1, 14, 16, 0, 0, time.UTC)
	if isTimeMatching(expr, nonMatchingTime) {
		t.Errorf("Did not expect time %v to match expression %s", nonMatchingTime, exprStr)
	}
}

// TestCronScheduler_StartStop tests starting and stopping the scheduler.
func TestCronScheduler_StartStop(t *testing.T) {
	scheduler := NewCronScheduler()
	scheduler.Start()
	if !scheduler.running {
		t.Errorf("Scheduler should be running after Start()")
	}
	scheduler.Stop()
	if scheduler.running {
		t.Errorf("Scheduler should not be running after Stop()")
	}
}

// TestCronScheduler_Concurrency tests if tasks are executed concurrently.
func TestCronScheduler_Concurrency(t *testing.T) {
	scheduler := NewCronScheduler()

	var counter int
	var mu sync.Mutex

	task := func() {
		mu.Lock()
		counter++
		mu.Unlock()
	}

	// Schedule the task to run at the same time
	now := time.Now()
	minute := (now.Minute() + 1) % 60
	cronExpr := fmt.Sprintf("%d %d * * *", minute, now.Hour())

	for i := 0; i < 5; i++ {
		err := scheduler.AddJob(cronExpr, task)
		if err != nil {
			t.Fatalf("Failed to add job: %v", err)
		}
	}

	scheduler.Start()
	defer scheduler.Stop()

	// Wait for tasks to be executed or timeout after 70 seconds
	time.Sleep(70 * time.Second)

	mu.Lock()
	if counter != 5 {
		t.Errorf("Expected counter to be 5, got %d", counter)
	}
	mu.Unlock()
}
