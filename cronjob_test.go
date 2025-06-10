package cronjob

import (
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
		{"* * * * * *", true},
		{"0 */15 * * * *", true},
		{"0 0 12 * * Mon-Fri", true},
		{"0 0 0 1 1 *", true},
		{"0 0 0 1 Jan *", true},
		{"0 0 0 * * Sun", true},
		{"0 0 0 * * 0", true},          // Sunday as 0
		{"0 0 0 * * 6", true},          // Saturday as 6
		{"60 * * * * *", false},        // Invalid second
		{"* 60 * * * *", false},        // Invalid minute
		{"* * 24 * * *", false},        // Invalid hour
		{"* * * 32 * *", false},        // Invalid day of month
		{"* * * * 13 *", false},        // Invalid month
		{"* * * * * SunFunday", false}, // Invalid day of week
		{"* * * * *", false},           // Too few fields (5 instead of 6)
		{"* * *", false},               // Too few fields
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

	err := scheduler.AddJob("* * * * * *", func() {})
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

	_ = scheduler.AddJob("* * * * * *", func() {})
	_ = scheduler.AddJob("0 */5 * * * *", func() {})

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
	executed := false

	// Job that runs every second
	err := scheduler.AddJob("* * * * * *", func() {
		executed = true
	})
	if err != nil {
		t.Fatalf("Failed to add job: %v", err)
	}

	scheduler.Start()
	defer scheduler.Stop()

	// Wait for the job to execute
	time.Sleep(2 * time.Second)

	if !executed {
		t.Errorf("Expected job to be executed")
	}
}

// TestIsTimeMatching tests the time matching function directly.
func TestIsTimeMatching(t *testing.T) {
	exprStr := "30 15 14 1 1 *" // At 14:15:30 on the 1st of January
	expr, err := ParseCronExpression(exprStr)
	if err != nil {
		t.Fatalf("Failed to parse cron expression: %v", err)
	}

	matchingTime := time.Date(2023, time.January, 1, 14, 15, 30, 0, time.UTC)
	if !isTimeMatching(expr, matchingTime) {
		t.Errorf("Expected time %v to match expression %s", matchingTime, exprStr)
	}

	nonMatchingTime := time.Date(2023, time.January, 1, 14, 15, 31, 0, time.UTC)
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

	// Add multiple jobs concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := scheduler.AddJob("* * * * * *", func() {
				// Simple task
			})
			if err != nil {
				t.Errorf("Failed to add job: %v", err)
			}
		}(i)
	}

	wg.Wait()

	if len(scheduler.Jobs) != 10 {
		t.Errorf("Expected 10 jobs, got %d", len(scheduler.Jobs))
	}
}
