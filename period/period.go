package period

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Period describes a period element for which you can test if a value is included
type Period interface {
	fmt.Stringer
	In(v int) bool
}

// PeriodeValue defines a fixed value
type PeriodeValue struct {
	Value int
}

func (p PeriodeValue) In(v int) bool {
	return p.Value == v
}

func (p PeriodeValue) String() string {
	return fmt.Sprintf("%d", p.Value)
}

func (p PeriodeValue) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// PeriodRange defines a range of values
type PeriodRange struct {
	Min int
	Max int
}

// In check if a value is in included in range
func (p PeriodRange) In(v int) bool {
	return v >= p.Min && v <= p.Max
}

func (p PeriodRange) String() string {
	return fmt.Sprintf("%d-%d", p.Min, p.Max)
}

func (p PeriodRange) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// PeriodDefinition  Defines criteria to match a ref time
// Like a Cron period definition.
type PeriodDefinition struct {
	Minutes []Period
	Hours   []Period
	Days    []Period
	Months  []Period
	Dows    []Period
}

// In checks if ref time matches the period criteria
func (p *PeriodDefinition) In(ref time.Time) bool {
	inside := true
	if inside && len(p.Minutes) > 0 {
		inside = inside && inValue(p.Minutes, ref.Minute())
	}
	if inside && len(p.Hours) > 0 {
		inside = inside && inValue(p.Hours, ref.Hour())
	}
	if inside && len(p.Days) > 0 {
		inside = inside && inValue(p.Days, ref.Day())
	}
	if inside && len(p.Dows) > 0 {
		inside = inside && inValue(p.Dows, int(ref.Weekday()))
	}
	if inside && len(p.Months) > 0 {
		inside = inside && inValue(p.Months, int(ref.Month()))
	}
	return inside
}

func (p *PeriodDefinition) String() string {
	min := periodsToText(p.Minutes)
	h := periodsToText(p.Hours)
	days := periodsToText(p.Days)
	months := periodsToText(p.Months)
	dows := periodsToText(p.Dows)
	return fmt.Sprintf("%s %s %s %s %s", min, h, days, months, dows)
}

func (p *PeriodDefinition) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

func periodsToText(periods []Period) string {
	if len(periods) == 0 {
		return "*"
	}
	b := strings.Builder{}
	for i, p := range periods {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(p.String())
	}
	return b.String()
}

func inValue(pp []Period, v int) bool {
	for _, p := range pp {
		if p.In(v) {
			return true
		}
	}
	return false
}

// MatchPeriods check if ref time matches any PeriodDefinition
func MatchPeriods(pp []PeriodDefinition, ref time.Time) bool {
	for _, p := range pp {
		if p.In(ref) {
			return true
		}
	}
	return false
}

// ParsePeriod parses string into Period element, handling fixed value or range min-max
func ParsePeriod(spec string, min int, max int) ([]Period, error) {
	values := make([]Period, 0, 2)

	parts := strings.Split(spec, ",")

	for _, p := range parts {
		var v Period
		if strings.Contains(p, "-") {
			vv := strings.Split(p, "-")
			v1, err := strconv.Atoi(vv[0])
			if err != nil {
				return values, err
			}
			v2, err := strconv.Atoi(vv[1])
			if err != nil {
				return values, err
			}
			v = PeriodRange{Min: v1, Max: v2}

		} else {
			i, err := strconv.Atoi(p)
			if err != nil {
				return values, err
			}
			v = PeriodeValue{Value: i}
		}
		values = append(values, v)
	}
	return values, nil
}
