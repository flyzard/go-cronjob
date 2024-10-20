// Package main is the entry point of the application.
package main

import (
	"fmt"
	"time"

	cronjob "github.com/flyzard/go-cronjob" // Replace with your actual module path
)

func main() {
	// Create a new CronScheduler instance
	scheduler := cronjob.NewCronScheduler()

	// Add a job that runs every minute
	err := scheduler.AddJob("* * * * *", func() {
		fmt.Println("Task 1: Runs every minute -", time.Now().Format(time.RFC1123))
	})
	if err != nil {
		fmt.Println("Error adding Task 1:", err)
		return
	}

	// Add a job that runs at 9 AM every Monday
	err = scheduler.AddJob("0 9 * * Mon", func() {
		fmt.Println("Task 2: Runs at 9 AM every Monday -", time.Now().Format(time.RFC1123))
	})
	if err != nil {
		fmt.Println("Error adding Task 2:", err)
		return
	}

	// Add a job that runs every 15 minutes
	err = scheduler.AddJob("*/15 * * * *", func() {
		fmt.Println("Task 3: Runs every 15 minutes -", time.Now().Format(time.RFC1123))
	})
	if err != nil {
		fmt.Println("Error adding Task 3:", err)
		return
	}

	// Add a job that runs at midnight on the first day of every month
	err = scheduler.AddJob("0 0 1 * *", func() {
		fmt.Println("Task 4: Runs at midnight on the first day of every month -", time.Now().Format(time.RFC1123))
	})
	if err != nil {
		fmt.Println("Error adding Task 4:", err)
		return
	}

	// Add a job with named month and day of week
	err = scheduler.AddJob("30 14 15 Jan-Mar Fri", func() {
		fmt.Println("Task 5: Runs at 14:30 on the 15th day of Jan, Feb, Mar and every Friday -", time.Now().Format(time.RFC1123))
	})
	if err != nil {
		fmt.Println("Error adding Task 5:", err)
		return
	}

	// Add a job that will panic to demonstrate panic handling
	err = scheduler.AddJob("2 * * * *", func() {
		fmt.Println("Task 6: This task will panic -", time.Now().Format(time.RFC1123))
		panic("intentional panic for testing")
	})
	if err != nil {
		fmt.Println("Error adding Task 6:", err)
		return
	}

	// Start the scheduler
	scheduler.Start()
	fmt.Println("CronScheduler started...")

	// Optionally, list all scheduled jobs
	fmt.Println("Scheduled Jobs:")
	jobs := scheduler.ListJobs()
	for _, job := range jobs {
		fmt.Println(job)
	}

	// Keep the main function running indefinitely
	// You can replace this with more sophisticated control if needed
	select {}
}
