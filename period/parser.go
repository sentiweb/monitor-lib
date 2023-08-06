package period

import "fmt"

// PeriodFields describe names of fields for each field, only used for error output
type PeriodFields struct {
	Minutes string
	Hours   string
	Days    string
	Months  string
	Dows    string
}

type PeriodParser struct {
	fields PeriodFields
}

func NewPeriodParser(fields PeriodFields) *PeriodParser {
	if fields.Minutes == "" {
		fields.Minutes = "min"
	}
	if fields.Hours == "" {
		fields.Hours = "hour"
	}
	if fields.Days == "" {
		fields.Days = "day"
	}
	if fields.Months == "" {
		fields.Months = "month"
	}
	if fields.Dows == "" {
		fields.Dows = "dow"
	}
	return &PeriodParser{fields: fields}
}

func Err(field string, msg error) error {
	return fmt.Errorf("%s : %s", field, msg)
}

func (parser *PeriodParser) Parse(Minutes string, Hours string, Days string, Months string, Dows string) (PeriodDefinition, error) {
	p := PeriodDefinition{}
	var (
		v   []Period
		err error
	)
	if Months != "" {
		v, err = ParsePeriod(Months, 1, 12)
		if err != nil {
			return p, Err(parser.fields.Months, err)

		}
		p.Months = v
	}
	if Minutes != "" {
		v, err = ParsePeriod(Minutes, 0, 59)
		if err != nil {
			return p, Err(parser.fields.Minutes, err)
		}
		p.Minutes = v
	}

	if Hours != "" {
		v, err = ParsePeriod(Hours, 0, 23)
		if err != nil {
			return p, Err(parser.fields.Hours, err)
		}
		p.Hours = v
	}
	if Days != "" {
		v, err = ParsePeriod(Days, 1, 31)
		if err != nil {
			return p, Err(parser.fields.Days, err)
		}
		p.Days = v
	}

	if Dows != "" {
		v, err = ParsePeriod(Dows, 0, 6)
		if err != nil {
			return p, Err(parser.fields.Dows, err)
		}
		p.Dows = v
	}
	return p, nil
}
