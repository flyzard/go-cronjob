// Package cronjob implements a cron expression parser.
package cronjob

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CronExpression represents a cron expression.
type CronExpression struct {
	Minutes    []int
	Hours      []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
}

var monthNameToNumber = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}

var dayNameToNumber = map[string]int{
	"Sun": 0,
	"Mon": 1,
	"Tue": 2,
	"Wed": 3,
	"Thu": 4,
	"Fri": 5,
	"Sat": 6,
}

// ParseCronExpression parses a cron expression and returns a CronExpression object.
func ParseCronExpression(expr string) (*CronExpression, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: %s", expr)
	}

	minutes, err := parseField(fields[0], 0, 59, nil)
	if err != nil {
		return nil, err
	}

	hours, err := parseField(fields[1], 0, 23, nil)
	if err != nil {
		return nil, err
	}

	dayOfMonth, err := parseField(fields[2], 1, 31, nil)
	if err != nil {
		return nil, err
	}

	month, err := parseField(fields[3], 1, 12, monthNameToNumber)
	if err != nil {
		return nil, err
	}

	dayOfWeek, err := parseField(fields[4], 0, 6, dayNameToNumber)
	if err != nil {
		return nil, err
	}

	return &CronExpression{
		Minutes:    minutes,
		Hours:      hours,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
	}, nil
}

func parseField(field string, min, max int, nameToNumber map[string]int) ([]int, error) {
	if field == "*" {
		var values []int
		for i := min; i <= max; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	if strings.Contains(field, "/") {
		return parseStepField(field, min, max, nameToNumber)
	}

	var values []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeValues, err := parseRange(part, min, max, nameToNumber)
			if err != nil {
				return nil, err
			}
			values = append(values, rangeValues...)
		} else {
			num, err := parseValue(part, min, max, nameToNumber)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		}
	}
	return values, nil
}

func parseValue(part string, min, max int, nameToNumber map[string]int) (int, error) {

	caser := cases.Title(language.Und)

	// Check if the part is a named value
	if num, ok := nameToNumber[caser.String(strings.ToLower(part))]; ok {
		return num, nil
	}

	// Try parsing as integer
	num, err := strconv.Atoi(part)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %s", part)
	}
	if num < min || num > max {
		return 0, fmt.Errorf("value out of range: %d", num)
	}
	return num, nil
}

func parseStepField(field string, min, max int, nameToNumber map[string]int) ([]int, error) {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid step field: %s", field)
	}

	baseValues, err := parseField(parts[0], min, max, nameToNumber)
	if err != nil {
		return nil, err
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return nil, fmt.Errorf("invalid step: %s", parts[1])
	}

	var values []int
	for _, val := range baseValues {
		if (val-min)%step == 0 {
			values = append(values, val)
		}
	}
	return values, nil
}

func parseRange(part string, min, max int, nameToNumber map[string]int) ([]int, error) {
	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", part)
	}

	start, err := parseValue(rangeParts[0], min, max, nameToNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid range start: %s", rangeParts[0])
	}

	end, err := parseValue(rangeParts[1], min, max, nameToNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid range end: %s", rangeParts[1])
	}

	if start > end {
		// For ranges like Fri-Mon, which should wrap around
		var values []int
		for i := start; i <= max; i++ {
			values = append(values, i)
		}
		for i := min; i <= end; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	var values []int
	for i := start; i <= end; i++ {
		values = append(values, i)
	}
	return values, nil
}
