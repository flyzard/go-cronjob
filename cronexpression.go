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
	Seconds    []int
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
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid cron expression: expected 6 fields (seconds minutes hours day month weekday), got %d", len(fields))
	}

	seconds, err := parseField(fields[0], 0, 59, nil)
	if err != nil {
		return nil, err
	}

	minutes, err := parseField(fields[1], 0, 59, nil)
	if err != nil {
		return nil, err
	}

	hours, err := parseField(fields[2], 0, 23, nil)
	if err != nil {
		return nil, err
	}

	dayOfMonth, err := parseField(fields[3], 1, 31, nil)
	if err != nil {
		return nil, err
	}

	month, err := parseField(fields[4], 1, 12, monthNameToNumber)
	if err != nil {
		return nil, err
	}

	dayOfWeek, err := parseField(fields[5], 0, 6, dayNameToNumber)
	if err != nil {
		return nil, err
	}

	return &CronExpression{
		Seconds:    seconds,
		Minutes:    minutes,
		Hours:      hours,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
	}, nil
}

func parseField(field string, minVal, maxVal int, nameToNumber map[string]int) ([]int, error) {
	if field == "*" {
		var values []int
		for i := minVal; i <= maxVal; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	if strings.Contains(field, "/") {
		return parseStepField(field, minVal, maxVal, nameToNumber)
	}

	var values []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeValues, err := parseRange(part, minVal, maxVal, nameToNumber)
			if err != nil {
				return nil, err
			}
			values = append(values, rangeValues...)
		} else {
			num, err := parseValue(part, minVal, maxVal, nameToNumber)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		}
	}
	return values, nil
}

func parseValue(part string, minVal, maxVal int, nameToNumber map[string]int) (int, error) {

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
	if num < minVal || num > maxVal {
		return 0, fmt.Errorf("value out of range: %d", num)
	}
	return num, nil
}

func parseStepField(field string, minVal, maxVal int, nameToNumber map[string]int) ([]int, error) {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid step field: %s", field)
	}

	baseValues, err := parseField(parts[0], minVal, maxVal, nameToNumber)
	if err != nil {
		return nil, err
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return nil, fmt.Errorf("invalid step: %s", parts[1])
	}

	var values []int
	for _, val := range baseValues {
		if (val-minVal)%step == 0 {
			values = append(values, val)
		}
	}
	return values, nil
}

func parseRange(part string, minVal, maxVal int, nameToNumber map[string]int) ([]int, error) {
	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", part)
	}

	start, err := parseValue(rangeParts[0], minVal, maxVal, nameToNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid range start: %s", rangeParts[0])
	}

	end, err := parseValue(rangeParts[1], minVal, maxVal, nameToNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid range end: %s", rangeParts[1])
	}

	if start > end {
		// For ranges like Fri-Mon, which should wrap around
		var values []int
		for i := start; i <= maxVal; i++ {
			values = append(values, i)
		}
		for i := minVal; i <= end; i++ {
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
