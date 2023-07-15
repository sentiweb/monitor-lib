package period

import (
	"fmt"
	"testing"
	"time"
)

func assertBool(t *testing.T, value bool, expect bool, msg string) {
	if value != expect {
		t.Errorf("Expected %t, got %t for %s", expect, value, msg)
	}
}

func assertIn(t *testing.T, p Period, value int, expect bool) {
	assertBool(t, p.In(value), expect, fmt.Sprintf("In(%d)", value))
}

func TestPeriod(t *testing.T) {

	p := PeriodeValue{Value: 1}

	assertIn(t, p, 1, true)
	assertIn(t, p, 2, false)
	assertIn(t, p, 5, false)
	assertIn(t, p, 0, false)

	r := PeriodRange{Min: 1, Max: 5}

	assertIn(t, r, 0, false)
	assertIn(t, r, 1, true)
	assertIn(t, r, 2, true)
	assertIn(t, r, 5, true)
	assertIn(t, r, 7, false)

}

func periods(pp ...Period) []Period {
	r := make([]Period, 0, len(pp))
	for _, p := range pp {
		r = append(r, p)
	}
	return r
}

func assertPeriod(t *testing.T, ref time.Time, p PeriodDefinition , expect bool) {
	r := p.In(ref)
	if r != expect {
		t.Errorf("Time %v <=> %+v, expected %t got %t", ref, p, expect, r)
	}
}

func date(year int, m time.Month, d int, h int, mn int) time.Time {
	return time.Date(year, m, d, h, mn, 0, 0, time.UTC)
}

func TestHours(t *testing.T) {
	// Hours
	i1 := PeriodDefinition {Hours: periods(PeriodRange{Min: 10, Max: 12})}

	hh := map[int]bool{
		1:  false,
		5:  false,
		10: true,
		11: true,
		12: true,
		13: false,
		22: false,
	}

	for h, expect := range hh {
		t1 := date(2010, time.January, 10, h, 0)
		assertPeriod(t, t1, i1, expect)

		t2 := date(2010, time.January, 10, h, 22)
		assertPeriod(t, t2, i1, expect)
	}
}

func TestMinutes(t *testing.T) {
	// Minutes
	i2 := PeriodDefinition {Minutes: periods(PeriodRange{Min: 20, Max: 30})}

	mm := map[int]bool{
		1:  false,
		5:  false,
		10: false,
		15: false,
		20: true,
		22: true,
		25: true,
		30: true,
		31: false,
		58: false,
	}

	for m, expect := range mm {
		t1 := date(2010, time.January, 10, 10, m)
		assertPeriod(t, t1, i2, expect)

		t2 := date(2010, time.January, 10, 22, m)
		assertPeriod(t, t2, i2, expect)
	}
}

func TestDow(t *testing.T) {
	// Minutes
	i3 := PeriodDefinition {Dows: periods(PeriodRange{Min: 4, Max: 5})}

	mm := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: true,
		5: true,
		6: false,
		7: false,
	}

	for m, expect := range mm {
		// 2021 Nov 1 = Monday
		t1 := date(2021, time.November, m, 10, 20)
		assertPeriod(t, t1, i3, expect)

		t2 := date(2010, time.November, m, 22, 55)
		assertPeriod(t, t2, i3, expect)
	}
}

func TestIgnore(t *testing.T) {

	TestHours(t)
	TestMinutes(t)
	TestDow(t)

}
