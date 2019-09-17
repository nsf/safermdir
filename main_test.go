package main

import (
	"fmt"
	"testing"
)

func TestPathCheck(t *testing.T) {
	type testCase struct {
		a       string
		b       string
		isError bool
	}
	var cases = []testCase{
		{"..", "b", true},
		{"", "b", true},
		{".", "b", true},
		{" . ", "b", true},
		{"/", "b", true},
		{"/ ", "b", true},
		{" /", "b", true},
		{" / ", "b", true},
		{"a", "b", true},
		{"/a", "b", false},

		{"/a", "..", true},
		{"/a", "", true},
		{"/a", ".", true},
		{"/a", " .", true},
		{"/a", ". ", true},
		{"/a", " . ", true},
		{"/a", "/", true},
		{"/a", " /", true},
		{"/a", "/ ", true},
		{"/a", " / ", true},
		{"/a", "/b", true},
		{"/a", "b", false},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := pathCheck(c.a, c.b)
			hasErr := err != nil
			if hasErr != c.isError {
				if hasErr {
					t.Errorf("unexpected error: %s", err)
				} else {
					t.Errorf("expected error, but got nothing")
				}
			}
		})
	}
}
