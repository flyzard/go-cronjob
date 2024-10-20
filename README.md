# Go-Cronjob

![Go](https://img.shields.io/badge/Go-v1.18-blue.svg)
![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/flyzard/go-cronjob.svg)](https://pkg.go.dev/github.com/flyzard/go-cronjob)

**Go-Cronjob** is a lightweight and efficient cron scheduler for Go applications. It allows you to schedule and manage recurring tasks using standard cron expressions, enabling automation of repetitive operations with ease.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Example](#basic-example)
  - [Advanced Usage](#advanced-usage)
- [API Reference](#api-reference)
  - [CronScheduler](#cronscheduler)
  - [CronExpression](#cronexpression)
- [Cron Expression Format](#cron-expression-format)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- **Standard Cron Expressions:** Supports the familiar five-field cron syntax.
- **Concurrency Control:** Executes scheduled tasks concurrently without blocking the scheduler.
- **Panic Handling:** Gracefully handles panics within tasks to ensure scheduler stability.
- **Job Management:** Easily add, remove, and list scheduled jobs.
- **Thread-Safe:** Designed with concurrency in mind, ensuring safe operations across multiple goroutines.
- **Extensible:** Allows for future enhancements like persistent storage, web interfaces, and more.

## Installation

To install the `go-cronjob` package, use the `go get` command:

```bash
go get github.com/flyzard/go-cronjob@v1.0.0
```

Ensure that your project is using Go modules. If not, initialize a new module:

```bash
go mod init your_module_name
```

## Usage

### Basic Example

Here's a simple example demonstrating how to use `go-cronjob` to schedule tasks:

```go
package main

import (
    "log"
    "time"

    "github.com/flyzard/go-cronjob"
)

func main() {
    // Create a new CronScheduler instance
    scheduler := cronjob.NewCronScheduler()

    // Add a job that runs every minute
    err := scheduler.AddJob("* * * * *", func() {
        log.Println("Task: Runs every minute -", time.Now())
    })
    if err != nil {
        log.Fatal(err)
    }

    // Start the scheduler
    scheduler.Start()
    defer scheduler.Stop()

    log.Println("CronScheduler started...")

    // Keep the application running
    select {}
}
```

**Expected Output:**

```
2024/04/27 10:00:00 main.go:14: CronScheduler started...
2024/04/27 10:00:00 main.go:8: Task: Runs every minute - Sat, 27 Apr 2024 10:00:00 UTC
...
```

### Advanced Usage

You can schedule multiple tasks with different cron expressions:

```go
package main

import (
    "log"
    "time"

    "github.com/flyzard/go-cronjob"
)

func main() {
    scheduler := cronjob.NewCronScheduler()

    // Task 1: Runs every minute
    err := scheduler.AddJob("* * * * *", func() {
        log.Println("Task 1: Every minute -", time.Now())
    })
    if err != nil {
        log.Fatal(err)
    }

    // Task 2: Runs at 9 AM every Monday
    err = scheduler.AddJob("0 9 * * Mon", func() {
        log.Println("Task 2: 9 AM every Monday -", time.Now())
    })
    if err != nil {
        log.Fatal(err)
    }

    // Task 3: Runs every 15 minutes
    err = scheduler.AddJob("*/15 * * * *", func() {
        log.Println("Task 3: Every 15 minutes -", time.Now())
    })
    if err != nil {
        log.Fatal(err)
    }

    // Start the scheduler
    scheduler.Start()
    defer scheduler.Stop()

    log.Println("CronScheduler started...")

    // List all scheduled jobs
    jobs := scheduler.ListJobs()
    for _, job := range jobs {
        log.Println(job)
    }

    // Keep the application running
    select {}
}
```

**Output:**

```
2024/04/27 09:00:00 main.go:20: CronScheduler started...
2024/04/27 09:00:00 main.go:24: Job 0: Schedule {Minutes:[0 1 2 ... 59] Hours:[0 1 ... 23], DayOfMonth:[1 2 ... 31], Month:[1 2 ... 12], DayOfWeek:[0 1 2 3 4 5 6]}
...
2024/04/27 09:00:00 main.go:17: Task 1: Every minute - Sat, 27 Apr 2024 09:00:00 UTC
2024/04/27 09:00:00 main.go:21: Task 3: Every 15 minutes - Sat, 27 Apr 2024 09:00:00 UTC
```

## API Reference

### `CronScheduler`

The `CronScheduler` struct manages the scheduling and execution of cron jobs.

#### `NewCronScheduler() *CronScheduler`

Creates and returns a new instance of `CronScheduler`.

```go
func NewCronScheduler() *CronScheduler
```

#### `AddJob(expr string, task func()) error`

Adds a new job to the scheduler with the specified cron expression and task function.

- **Parameters:**
  - `expr`: A string representing the cron expression.
  - `task`: A function to execute when the cron expression matches.

- **Returns:**
  - `error`: An error if the cron expression is invalid or the job cannot be added.

```go
func (c *CronScheduler) AddJob(expr string, task func()) error
```

#### `RemoveJob(index int) error`

Removes the job at the specified index from the scheduler.

- **Parameters:**
  - `index`: The index of the job to remove.

- **Returns:**
  - `error`: An error if the index is out of range.

```go
func (c *CronScheduler) RemoveJob(index int) error
```

#### `ListJobs() []string`

Returns a list of all scheduled jobs in the scheduler.

- **Returns:**
  - `[]string`: A slice of strings describing each job.

```go
func (c *CronScheduler) ListJobs() []string
```

#### `Start()`

Starts the cron scheduler, enabling it to begin executing scheduled jobs.

```go
func (c *CronScheduler) Start()
```

#### `Stop()`

Stops the cron scheduler gracefully, ensuring that no jobs are left running.

```go
func (c *CronScheduler) Stop()
```

### `CronExpression`

The `CronExpression` struct represents a parsed cron expression.

#### Fields:

- `Minutes []int`: Allowed minutes (0-59).
- `Hours []int`: Allowed hours (0-23).
- `DayOfMonth []int`: Allowed days of the month (1-31).
- `Month []int`: Allowed months (1-12).
- `DayOfWeek []int`: Allowed days of the week (0-6, where 0 is Sunday).

```go
type CronExpression struct {
    Minutes    []int
    Hours      []int
    DayOfMonth []int
    Month      []int
    DayOfWeek  []int
}
```

#### `ParseCronExpression(expr string) (*CronExpression, error)`

Parses a cron expression string and returns a `CronExpression` object.

- **Parameters:**
  - `expr`: A string representing the cron expression.

- **Returns:**
  - `*CronExpression`: The parsed cron expression.
  - `error`: An error if the cron expression is invalid.

```go
func ParseCronExpression(expr string) (*CronExpression, error)
```

## Cron Expression Format

The cron expression follows the standard five-field format:

```
* * * * *
| | | | |
| | | | +----- Day of the Week (0 - 6) (Sunday=0)
| | | +------- Month (1 - 12)
| | +--------- Day of the Month (1 - 31)
| +----------- Hour (0 - 23)
+------------- Minute (0 - 59)
```

### Supported Syntax:

- **Asterisk (`*`):** Represents all possible values for a field.
- **Comma (`,`):** Specifies a list of values.
- **Dash (`-`):** Defines a range of values.
- **Slash (`/`):** Indicates step values.

### Examples:

- `* * * * *`: Every minute.
- `0 9 * * Mon`: At 9 AM every Monday.
- `*/15 * * * *`: Every 15 minutes.
- `30 14 15 Jan-Mar Fri`: At 14:30 on the 15th day of January through March and every Friday.

## Testing

The package includes comprehensive unit tests to ensure reliability and correctness.

### Running Tests

To execute the tests, navigate to the project directory and run:

```bash
go test ./...
```

### Test Coverage

The tests cover various scenarios, including:

- Parsing valid and invalid cron expressions.
- Scheduling and executing tasks.
- Handling edge cases like invalid ranges and steps.
- Ensuring panic handling within tasks.

## Contributing

Contributions are welcome! To contribute to the `go-cronjob` project, please follow these steps:

1. **Fork the Repository:**

   Click the "Fork" button at the top-right corner of the repository page to create a personal copy.

2. **Clone Your Fork:**

   ```bash
   git clone https://github.com/yourusername/go-cronjob.git
   cd go-cronjob
   ```

3. **Create a New Branch:**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make Your Changes:**

   Implement your feature or fix in the appropriate files.

5. **Run Tests:**

   Ensure all tests pass:

   ```bash
   go test ./...
   ```

6. **Commit Your Changes:**

   ```bash
   git add .
   git commit -m "Add your commit message"
   ```

7. **Push to Your Fork:**

   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create a Pull Request:**

   Navigate to your fork on GitHub and click "Compare & pull request." Provide a clear description of your changes.

### Guidelines

- **Code Quality:** Follow Go's best practices and maintain consistent coding standards.
- **Documentation:** Update or add documentation for any new features or changes.
- **Testing:** Add tests for new functionalities to ensure robustness.
- **Respect Licensing:** Ensure that all contributions comply with the project's MIT License.

## License

This project is licensed under the [MIT License](LICENSE).

```
MIT License

Copyright (c) 2024 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

[...rest of the MIT License...]
```

## Acknowledgements

- Inspired by traditional cron systems and the need for a simple, reliable scheduler in Go.
- Utilizes the [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) package for text manipulation.

---

## Getting Started

To get started with `go-cronjob`, refer to the [Usage](#usage) section above. For more detailed examples and API usage, visit the [Go Documentation](https://pkg.go.dev/github.com/flyzard/go-cronjob).

Feel free to open issues or submit pull requests for any enhancements or bug fixes. Happy scheduling!

---