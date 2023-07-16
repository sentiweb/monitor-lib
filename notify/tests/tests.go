package tests

import (
	"fmt"
	"time"
)

type MockService struct {
	Name string
}

func (m *MockService) String() string {
	return m.Name
}

type MockNotification struct {
	StatusValue   string
	LabelValue    string
	ServiceValue  *MockService
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

func (m *MockNotification) Service() fmt.Stringer {
	return m.ServiceValue
}

func (m *MockNotification) String() string {
	return m.LabelValue
}

func NewMockNotification(status string, label string, time time.Time) *MockNotification {
	return &MockNotification{StatusValue: status, LabelValue: label, ServiceValue: &MockService{Name: label}, FromTimeValue: time}
}
