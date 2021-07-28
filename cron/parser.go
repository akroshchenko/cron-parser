package cron

import (
	"fmt"
	"strconv"
	"strings"
)

const supportedFieldCharacters = "0123456789*-/,"

type bounds struct {
	min uint
	max uint
}

type field struct {
	name    string
	bounds  bounds
	ranges  []uint
	command string
}

func (f field) getFieldName() string {
	return f.name
}
func (f field) getFieldDetails() string {
	if f.name == "command" {
		return f.command
	}
	return fmt.Sprintf("%v", f.ranges)
}

type schedule [6]field

func (s schedule) String() string {
	var result []string
	for _, f := range s {
		result = append(result, fmt.Sprintf("%-14s %v", f.getFieldName(), f.getFieldDetails()))
	}
	return strings.Join(result, "\n")
}

func Parse(cronStr string) (*schedule, error) {

	defaultSchedule := schedule{
		field{
			name:   "minutes",
			bounds: bounds{0, 59},
		},
		field{
			name:   "hours",
			bounds: bounds{0, 23},
		},
		field{
			name:   "day of month",
			bounds: bounds{1, 31},
		},
		field{
			name:   "month",
			bounds: bounds{1, 12},
		},
		field{
			name:   "day of week",
			bounds: bounds{1, 7},
		},
		field{
			name: "command",
		},
	}

	if len(cronStr) == 0 {
		return nil, fmt.Errorf("Error: empty expression string")
	}

	normalizedString := strings.Join(strings.Fields(cronStr), " ")

	fields := strings.Split(normalizedString, " ")
	if len(fields) != len(defaultSchedule) {
		return nil, fmt.Errorf("Error: the number of provided field does not match the specification (should be provided %d fields). Provided: %s", len(defaultSchedule), cronStr)
	}

	for i := 0; i < len(defaultSchedule); i++ {
		if defaultSchedule[i].name == "command" {
			// Do not parse field if is command
			defaultSchedule[i].command = fields[i]
			continue
		}
		r, err := getRange(fields[i], defaultSchedule[i].bounds)
		if err != nil {
			return nil, fmt.Errorf("Error to get range from field '%s': %v", fields[i], err)
		}
		defaultSchedule[i].ranges = r
	}

	return &defaultSchedule, nil

}

func getRange(f string, b bounds) ([]uint, error) {

	var r []uint
	for _, c := range f {
		if !strings.ContainsRune(supportedFieldCharacters, rune(c)) {
			return nil, fmt.Errorf("Not allowed character '%q' for field %s", c, f)
		}
	}

	if strings.ContainsRune(f, rune(',')) {
		sequence := strings.Split(f, ",")
		r = make([]uint, len(sequence))
		for indx, numb := range sequence {
			number, err := strconv.Atoi(numb)
			if err != nil {
				return nil, fmt.Errorf("Cannot convert %s to integer: %s", numb, err)
			}
			r[indx] = uint(number)
		}
		return r, nil
	}

	if strings.ContainsRune(f, rune('-')) {

		step := 1

		boarders := strings.Split(f, "-")
		if len(boarders) != 2 {
			return nil, fmt.Errorf("Wrong syntax of using '-', should be 'n-n' where n is a number: %s", f)
		}
		min, err := strconv.Atoi(boarders[0])
		if err != nil {
			return nil, fmt.Errorf("Cannot parse the min value of the %s", f)
		}

		max, err := strconv.Atoi(boarders[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot parse the max value of the %s", f)
		}
		if min > max {
			return nil, fmt.Errorf("Error: value on the left of '-' should be lover then on the right: %v", f)
		}

		if uint(min) < b.min || uint(min) > b.max || uint(max) < b.min || uint(max) > b.max {
			return nil, fmt.Errorf("wrong field range: %s", f)
		}

		r = make([]uint, ((max-min)/step)+1)
		for i := 0; i < len(r); i += step {
			r[i] = uint(min + i)
		}
		return r, nil
	}

	if strings.ContainsRune(f, rune('/')) {
		boarders := strings.Split(f, "/")
		if len(boarders) != 2 {
			return nil, fmt.Errorf("Wrong syntax of using '/', should be 'n/n' where n is a number(first n could be '*'): %s", f)
		}
		var start, step uint
		if boarders[0] == "*" {
			start = 0
		} else {
			parsedV, err := strconv.Atoi(boarders[0])
			if err != nil {
				return nil, fmt.Errorf("Cannot parse %s: %s", f, err)
			}
			start = uint(parsedV)
		}
		parsedV, err := strconv.Atoi(boarders[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot parse %s: %s", f, err)
		}
		step = uint(parsedV)

		r = make([]uint, (b.max-start)/step+1)
		current := start
		for i := 0; i < len(r); i++ {
			if current > b.max {
				break
			}
			r[i] = current
			current += step
		}
		return r, nil
	}

	if len(f) == 1 && f == "*" {
		var step uint = 1
		r = make([]uint, (b.max-b.min)/step+1)
		current := b.min
		for i := 0; i < len(r); i++ {
			if current > b.max {
				break
			}
			r[i] = current
			current += step
		}
		return r, nil
	}

	n, err := strconv.Atoi(f)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse value %s as a number", f)
	}
	return []uint{uint(n)}, nil
}
