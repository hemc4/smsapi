package main

import "testing"

var smstests = []struct {
	from string
	to   string
	text string
	out  string
}{
	{"12345678", "87654321", "hello", ""},
	{"", "87654321", "hello", "from is missing"},
	{"12345678", "", "hello", "to is missing"},
	{"12345678", "87654321", "", "text is missing"},
	{"1234", "87654321", "hello", "from is invalid"},
	{"12341234123412341234", "87654321", "hello", "from is invalid"},
	{"12345678", "87654", "hello", "to is invalid"},
	{"12345678", "8765487654876548765487654", "hello", "to is invalid"},
	{"12345678", "87654321", "hellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutside", "text is invalid"},

}

func TestValidateFormData(t *testing.T) {

	for _, smstest := range smstests {
		s := ValidateFormData(smstest.from, smstest.to, smstest.text)
		if s != smstest.out {
			t.Errorf("validateformdata(%q,%q,%q) => %q want %q", smstest.from, smstest.to, smstest.text, s, smstest.out)
		}
	}
}

func TestIndex(t *testing.T) {
	t.Log("hello form testing ")
	if true != true {
		t.Error("Never show it")
	}
}
