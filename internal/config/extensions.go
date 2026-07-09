package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StringSlice is a custom type that can parse JSON arrays of strings
type StringSlice []string

func (slice *StringSlice) UnmarshalText(text []byte) error {
	value := string(text)
	if value == "" {
		*slice = []string{}
		return nil
	}

	var arr []string
	if err := json.Unmarshal(text, &arr); err == nil {
		*slice = arr
		return nil
	}

	// Fallback to comma-separated values
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	*slice = parts
	return nil
}

// IntSlice is a custom type that can parse JSON arrays of integers
type IntSlice []int

func (slice *IntSlice) UnmarshalText(text []byte) error {
	value := string(text)
	if value == "" {
		*slice = []int{}
		return nil
	}

	var arr []int
	if err := json.Unmarshal(text, &arr); err == nil {
		*slice = arr
		return nil
	}

	// Fallback to comma-separated values
	parts := strings.Split(value, ",")
	ints := make([]int, 0, len(parts))
	for _, part := range parts {
		num, ok := TryParseInteger(part)
		if !ok {
			return fmt.Errorf("unable to parse integer: %s", part)
		}
		ints = append(ints, num)
	}
	*slice = ints
	return nil
}

func TryParseInteger(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return num, true
}

// DynamicTime is a custom type that can parse multiple time formats
type DynamicTime time.Time

func (dt *DynamicTime) UnmarshalText(text []byte) error {
	value := string(text)
	if value == "" {
		*dt = DynamicTime(time.Time{})
		return nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			*dt = DynamicTime(t)
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", value)
}

func (dt DynamicTime) Time() time.Time {
	return time.Time(dt)
}
