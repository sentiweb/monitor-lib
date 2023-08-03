package sets

import (
	"fmt"
	"testing"
)

var PresentsSet = []string{"toto", "titi", "tutu", "tata", "112"}
var AbsentSet = []string{"12345", "", "myvalue"}

func AssertBool(t *testing.T, value bool, expected bool, msg string) {
	if value != expected {
		t.Errorf("Assertion '%s' faild, expected %t got %t", msg, expected, value)
	}
}

func Ensure(t *testing.T, s *Set[string], values []string, present bool) {
	var m string
	if present {
		m = "present"
	} else {
		m = "absent"
	}
	for _, v := range values {
		AssertBool(t, s.Has(v), present, fmt.Sprintf("Element '%s' should be %s", v, m))
	}
}

func TestSet(t *testing.T) {

	s := New[string](0)

	a := s.Add("toto")
	AssertBool(t, a, true, "Toto is added")

	AssertBool(t, s.Has("toto"), true, "Toto is present")

	Ensure(t, s, AbsentSet, false)
}

func TestAddAll(t *testing.T) {

	ff := PresentsSet

	s := New[string](0)

	a := s.AddAll(ff)
	AssertBool(t, a == len(ff), true, "Added element should be equal to array size")

	Ensure(t, s, PresentsSet, true)
	Ensure(t, s, AbsentSet, false)
}

func TestHasAny(t *testing.T) {

	ff := PresentsSet

	s := New[string](0)

	a := s.AddAll(ff)
	AssertBool(t, a == len(ff), true, "Added element should be equal to array size")

	AssertBool(t, s.HasAny(ff), true, "Should have any")
}
