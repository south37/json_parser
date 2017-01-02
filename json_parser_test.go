package json_parser

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		Input  string
		Expect bool
	}{
		{`{}`, true},
		{`{"str": 3}`, true},
		{`{"obj": {"nest": 0 }}`, true},
		{`[1, 2, 3]`, true},
		{`["1", 2, {"ok": 3}]`, true},
		{`{"array": [1, 2, 3]}`, true},
		{`"array"`, false},
	}
	for _, tc := range cases {
		parser := NewParser(tc.Input)
		err := parser.parse()
		if got := (err == nil); got != tc.Expect {
			t.Errorf("parse %s expected %t got %t, error is `%s`", tc.Input, tc.Expect, got, err)
		}
	}
}

func TestMatchString(t *testing.T) {
	cases := []struct {
		Input  string
		Expect bool
	}{
		{`""`, true},
		{`str`, false},
		{`st"r`, false},
		{`str"`, false},
		{`"str"`, true},
		{`"It's \"strong\" monster"`, true},
	}
	for _, tc := range cases {
		if got := matchString(tc.Input); got != tc.Expect {
			t.Errorf("%s: matchString expected %t, got %t", tc.Input, tc.Expect, got)
		}
	}
}

func TestMatchNumber(t *testing.T) {
	cases := []struct {
		Input  string
		Expect bool
	}{
		{`3`, true},
		{`-3`, true},
		{`0`, true},
		{`-0`, false},
		{`30`, true},
	}
	for _, tc := range cases {
		if got := matchNumber(tc.Input); got != tc.Expect {
			t.Errorf("%s: matchNumber expected %t, got %t", tc.Input, tc.Expect, got)
		}
	}
}
