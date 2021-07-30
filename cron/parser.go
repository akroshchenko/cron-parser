package cron

// TODO: document functions and structs
// TOOD: write tests
// TODO: add check for cases when day of month does get beyound the max days for this month
import (
	"fmt"
	"strconv"
	"strings"
)

const supportedFieldCharacters = "0123456789*-/,"

type bounds [2]uint // 0 - min, 1 - max

var (
	minuteBounds     = bounds{0, 59}
	hourBounds       = bounds{0, 23}
	dayOfMonthBounds = bounds{1, 31}
	monthBounds      = bounds{1, 12}
	dayOfWeekBounds  = bounds{1, 7}
)

func (b bounds) String() string {
	return fmt.Sprintf("[from %d to %d]", b[0], b[1])
}

func (b bounds) min() uint {
	return b[0]
}

func (b bounds) max() uint {
	return b[1]
}

func (b bounds) include(i uint) bool {
	return i >= b[0] && i <= b[1]
}

type field uint

const (
	minuteKey field = iota
	hourKey
	dayOfMonthKey
	monthKey
	dayOfWeekKey
	commandKey
)

// Q: is it ok to have type and methods of the same name?
func (f field) bounds() bounds {
	switch f {
	case minuteKey:
		return minuteBounds
	case hourKey:
		return hourBounds
	case dayOfMonthKey:
		return dayOfMonthBounds
	case monthKey:
		return monthBounds
	case dayOfWeekKey:
		return dayOfWeekBounds
	}
	return bounds{} // Q: what should I return here?
}

func (f field) String() string {
	switch f {
	case minuteKey:
		return "minute"
	case hourKey:
		return "hour"
	case dayOfMonthKey:
		return "day of month"
	case monthKey:
		return "month"
	case dayOfWeekKey:
		return "day of week"
	case commandKey:
		return "command"
	}

	return "" // Q: what should I return here?
}

type schedule struct {
	ranges  [5][]uint
	command string
}

func (s schedule) String() string {
	var result []string
	for k, r := range s.ranges {
		result = append(result, fmt.Sprintf("%-14s %v", field(k), strings.Trim(fmt.Sprint(r), "[]")))
	}
	result = append(result, fmt.Sprintf("%-14s %v", "command", s.command))
	return strings.Join(result, "\n")
}

func Parse(inputExpr string) (*schedule, error) {

	var scheduleSpec schedule

	if len(inputExpr) == 0 {
		return nil, fmt.Errorf("Error: empty expression string")
	}

	parsedFields := strings.Fields(inputExpr)

	if len(parsedFields) != (len(scheduleSpec.ranges) + 1) {
		return nil, fmt.Errorf("Error: the number of provided field does not match the specification (should be provided %d fields). Was provided: %s", len(scheduleSpec.ranges)+1, inputExpr)
	}

	for i := 0; i < len(parsedFields); i++ {
		if field(i) == commandKey {
			scheduleSpec.command = parsedFields[commandKey]
			continue
		}
		r, err := getRange(parsedFields[i], field(i).bounds())
		if err != nil {
			return nil, fmt.Errorf("Error: failed to get range from field '%s': %s with error: %s", field(i), parsedFields[i], err)
		}
		scheduleSpec.ranges[i] = make([]uint, 0, len(r))
		scheduleSpec.ranges[i] = r
	}

	return &scheduleSpec, nil

}

func getSequenceFromCommaEpr(expr string, b bounds) ([]uint, error) {
	var r []uint

	sequence := strings.Split(expr, ",")
	r = make([]uint, len(sequence))
	for indx, numb := range sequence {
		number, err := strconv.Atoi(numb)
		if err != nil {
			return nil, fmt.Errorf("Cannot convert %s to integer: %s", numb, err)
		}
		if !b.include(uint(number)) {
			return nil, fmt.Errorf("Number %d is not within bounds: %s", number, b)
		}
		r[indx] = uint(number)
	}
	return r, nil
}

func getSequenceFromHifenEpr(expr string, b bounds) ([]uint, error) {

	var r []uint

	boarders := strings.Split(expr, "-")
	if len(boarders) != 2 {
		return nil, fmt.Errorf("Wrong syntax of using '-', should be 'n-n' where n is a number: %s", expr)
	}
	min, err := strconv.Atoi(boarders[0])
	if err != nil {
		return nil, fmt.Errorf("Cannot parse the min value of the %s", expr)
	}

	if !b.include(uint(min)) {
		return nil, fmt.Errorf("Number %d is not within bounds: %s", min, b)
	}

	max, err := strconv.Atoi(boarders[1])
	if err != nil {
		return nil, fmt.Errorf("Cannot parse the max value of the %s", expr)
	}

	if !b.include(uint(max)) {
		return nil, fmt.Errorf("Number %d is not within bounds: %s", max, b)
	}

	if min > max {
		return nil, fmt.Errorf("Value on the left of '-' should be lover then on the right: %v", expr)
	}

	r = make([]uint, max-min+1)
	for i := 0; i < len(r); i++ {
		r[i] = uint(min + i)
	}
	return r, nil
}

func getSequenceFromSlashExpr(expr string, b bounds) ([]uint, error) {
	var r []uint
	boarders := strings.Split(expr, "/")
	if len(boarders) != 2 {
		return nil, fmt.Errorf("Wrong syntax of using '/', should be 'n/n' where n is a number(first n could be '*'): %s", expr)
	}
	var start, step uint
	if boarders[0] == "*" {
		start = b.min()
	} else {
		parsedStart, err := strconv.Atoi(boarders[0])
		if err != nil {
			return nil, fmt.Errorf("Cannot parse %s: %s", expr, err)
		}
		if !b.include(uint(parsedStart)) {
			return nil, fmt.Errorf("Number %d is not within bounds: %s", parsedStart, b)
		}
		start = uint(parsedStart)
	}
	parsedStep, err := strconv.Atoi(boarders[1])
	if err != nil {
		return nil, fmt.Errorf("Cannot parse %s: %s", expr, err)
	}
	if !b.include(uint(parsedStep)) {
		return nil, fmt.Errorf("Number %d is not within bounds: %s", parsedStep, b)
	}
	step = uint(parsedStep)
	if step == 0 {
		return nil, fmt.Errorf("The 'step' in the expression of format 'start/step' cannot be equeal to zero, got %s", expr)
	}

	r = make([]uint, (b.max()-start)/step+1)
	current := start
	for i := 0; i < len(r); i++ {
		if current > b.max() {
			break
		}
		r[i] = current
		current += step
	}
	return r, nil
}

func getRange(f string, b bounds) ([]uint, error) {

	var r []uint
	for _, c := range f {
		if !strings.ContainsRune(supportedFieldCharacters, rune(c)) {
			return nil, fmt.Errorf("Not allowed character '%q' for field %s", c, f)
		}
	}

	if strings.ContainsRune(f, rune(',')) {
		return getSequenceFromCommaEpr(f, b)
	}

	if strings.ContainsRune(f, rune('-')) {
		return getSequenceFromHifenEpr(f, b)
	}

	if strings.ContainsRune(f, rune('/')) {
		return getSequenceFromSlashExpr(f, b)
	}

	if len(f) == 1 && f == "*" {
		r = make([]uint, b.max()-b.min())
		current := b.min()
		for i := 0; i < len(r); i++ {
			if current > b.max() {
				break
			}
			r[i] = current
			current++
		}
		return r, nil
	}

	n, err := strconv.Atoi(f)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse value %s as a number", f)
	}
	if !b.include(uint(n)) {
		return nil, fmt.Errorf("Number %d is not within bounds: %s", n, b)
	}
	return []uint{uint(n)}, nil
}
