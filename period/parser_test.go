package period

import (
	"testing"
)

var fields = PeriodFields{}

func testParser(t *testing.T, ff [5]string, expect string) {

	p := NewPeriodParser(fields)

	d, err := p.Parse(ff[0], ff[1], ff[2], ff[3], ff[4])
	if err != nil {
		t.Error(err)
	}

	if d.String() != expect {
		t.Errorf("Bad result, got '%s', expect '%s'", d.String(), expect)
	}

}

func TestParser(t *testing.T) {

	testParser(t, [5]string{"1-3", "", "", "", ""}, "1-3 * * * *")
	testParser(t, [5]string{"", "5", "", "", ""}, "* 5 * * *")
	testParser(t, [5]string{"", "", "6", "", ""}, "* * 6 * *")
	testParser(t, [5]string{"", "", "", "7", ""}, "* * * 7 *")
	testParser(t, [5]string{"", "", "", "", "0"}, "* * * * 0")

}
