package tests

import (
	"time"
)

type MockNotification struct {
	StatusValue   string
	LabelValue    string
	ServiceValue  string
	FromTimeValue time.Time
}

func (m *MockNotification) Status() string {
	return m.StatusValue
}

func (m *MockNotification) Label() string {
	return m.LabelValue
}

func (m *MockNotification) FromTime() time.Time {
	return m.FromTimeValue
}

func (m *MockNotification) Tags() []string {
	return []string{}
}

func (m *MockNotification) ServiceName() string {
	return m.ServiceValue
}

func (m *MockNotification) String() string {
	return m.LabelValue
}

func NewMockNotification(service string, status string, label string, time time.Time) *MockNotification {
	return &MockNotification{StatusValue: status, LabelValue: label, ServiceValue: label, FromTimeValue: time}
}
