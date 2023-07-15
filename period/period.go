package period

import (
	"time"
	"strconv"
	"strings"
)

// Period describes a period element for which you can test if a value is included
type Period interface {
	In(v int) bool
}

// PeriodeValue defines a fixed value
type PeriodeValue struct {
	Value int
}

func (p PeriodeValue) In(v int) bool {
	return p.Value == v
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

// PeriodDefinition  Defines criteria to match a ref time
// Like a Cron period definition.
type PeriodDefinition  struct {
	Minutes []Period
	Hours   []Period
	Days    []Period
	Months  []Period
	Dows    []Period
}

// In checks if ref time matches the period criteria
func (p PeriodDefinition ) In(ref time.Time) bool {
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

func inValue(pp []Period, v int) bool {
	for _, p := range pp {
		if p.In(v) {
			return true
		}
	}
	return false
}

// MatchPeriods check if ref time matches any PeriodDefinition
func MatchPeriods(pp []PeriodDefinition , ref time.Time) bool {
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

