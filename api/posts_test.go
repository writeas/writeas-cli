package api

import "testing"

func TestTrimToLength(t *testing.T) {
	tt := []struct {
		Name       string
		Data       string
		Length     int
		ResultData string
		ResultIDX  int
	}{
		{
			"English, longer than trim length",
			"This is a string, let's truncate it.",
			12,
			"This is a",
			9,
		}, {
			"English, equal to length",
			"Some other string.",
			18,
			"Some other string.",
			-1,
		}, {
			"English, shorter than trim length",
			"I'm short!",
			20,
			"I'm short!",
			-1,
		}, {
			"Multi-byte, longer than trim length",
			"這是一個較長的廣東話。 有許多特性可以確保足夠長的輸出。",
			14,
			"這是一個較長的廣東話。",
			11,
		}, {
			"Multi-byte, equal to length",
			"這是一個簡短的廣東話。",
			11,
			"這是一個簡短的廣東話。",
			-1,
		}, {
			"Multi-byte, shorter than trim length",
			"我也很矮！ 有空間。",
			20,
			"我也很矮！ 有空間。",
			-1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			out, idx := trimToLength(tc.Data, tc.Length)
			if out != tc.ResultData {
				t.Errorf("Incorrect output, expecting \"%s\" but got \"%s\"", tc.ResultData, out)
			}

			if idx != tc.ResultIDX {
				t.Errorf("Incorrect last index, expected \"%d\" but got \"%d\"", tc.ResultIDX, idx)
			}
		})
	}
}

func TestGetExcerpt(t *testing.T) {
	tt := []struct {
		Name   string
		Data   string
		Result string
	}{
		{
			"Shorter than one line",
			"This is much less than 80 chars",
			"This is much less than 80 chars",
		}, {
			"Exact length, one line",
			"This will be only 80 chars. Maybe all the way to column 88, that will do it. ---",
			"This will be only 80 chars. Maybe all the way to column 88, that will do it. ---",
		}, {
			"Shorter than two lines",
			"This will be more than one line but shorter than two. It should break at the 80th or less character. Let's check it out.",
			"This will be more than one line but shorter than two. It should break at the\n 80th or less character. Let's check it out.",
		}, {
			"Exact length, two lines",
			"This should be the exact length for two lines. There should ideally be no trailing periods to indicate further text. However trimToLength breaks on word bounds.",
			"This should be the exact length for two lines. There should ideally be no\n trailing periods to indicate further text. However trimToLength breaks on word...",
		}, {
			"Longer than two lines",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque volutpat sagittis aliquet. Ut eu rutrum nisl. Proin molestie ante in dui vulputate dictum. Proin ac bibendum eros. Nulla porta congue tellus, sed vehicula sem bibendum eu. Donec vehicula erat viverra fermentum mattis. Integer volutpat.",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque volutpat\n sagittis aliquet. Ut eu rutrum nisl. Proin molestie ante in dui vulputate...",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			out := getExcerpt(tc.Data)
			if out != tc.Result {
				t.Errorf("Output does not match:\nexpected \"%s\"\nbut got \"%s\"", tc.Result, out)
			}
		})
	}
}
