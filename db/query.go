package db

type QueryOption struct{}

type IDQueryOption struct {
	IDs []int64
}

func NewIDQueryOption(ids []int64) *IDQueryOption {
	return &IDQueryOption{IDs: ids}
}

func (i *IDQueryOption) Values() []int64 {
	return i.IDs
}

func (i *IDQueryOption) IsEmpty() bool {
	return len(i.IDs) == 0
}

type IntQueryOption struct {
	value int64
}

func NewIntQueryOption(value int64) *IntQueryOption {
	return &IntQueryOption{value: value}
}

func (i *IntQueryOption) Value() int64 {
	return i.value
}

type RangeOperator uint8

const (
	TimeOpBefore RangeOperator = 1
	TimeOpAfter  RangeOperator = 2
	TimeOpRange  RangeOperator = 3
)

type RangeQueryOption struct {
	Op  RangeOperator
	Min int64 // Use for Before and After operators
	Max int64 // Only used for range Operator
}

func NewRangeQueryOption(op RangeOperator, min int64, max int64) *RangeQueryOption {
	return &RangeQueryOption{Op: op, Min: min, Max: max}
}

func NewTimeQueryOptionBefore(value int64) *RangeQueryOption {
	return &RangeQueryOption{Op: TimeOpBefore, Min: value}
}

func NewTimeQueryOptionAfter(value int64) *RangeQueryOption {
	return &RangeQueryOption{Op: TimeOpAfter, Min: value}
}

func NewTimeQueryOptionRange(min int64, max int64) *RangeQueryOption {
	return &RangeQueryOption{Min: min, Max: max, Op: TimeOpRange}
}
