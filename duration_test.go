package requests

import (
	"testing"
	"time"
)

const floatTimeDev = time.Duration(1) * time.Nanosecond

var ftosTable = []struct {
	n        float64
	expected time.Duration
}{
	{0.1, time.Duration(100) * time.Millisecond},
	{1.0, time.Duration(1) * time.Second},
	{10.0, time.Duration(10) * time.Second},
	{60.0, time.Duration(60) * time.Second},
}

func TestFtos(t *testing.T) {
	for _, tt := range ftosTable {
		result := Ftos(tt.n)
		if !(result >= tt.expected-floatTimeDev || result <= tt.expected-floatTimeDev) {
			logExpectedResult(result, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}

func TestStof(t *testing.T) {
	for _, tt := range ftosTable {
		result := Stof(tt.expected)
		if !(result >= tt.n-nano || result <= tt.n+nano) {
			logExpectedResult(result, tt.n)
			t.Error("Result isn't accurate.")
		}
	}
}

var hoursTable = []struct {
	n        Second
	expected float64
}{
	{3600.0, 1.0},
	{7200.0, 2.0},
	{10800.0, 3.0},
	{14400.0, 4.0},
	{18000.0, 5.0},
}

func TestHours(t *testing.T) {
	for _, tt := range hoursTable {
		s := tt.n
		if s.Hours() != tt.expected {
			logExpectedResult(s, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}

var minutesTable = []struct {
	n        Second
	expected float64
}{
	{30.0, 0.5},
	{60.0, 1.0},
	{90.0, 1.5},
	{120.0, 2.0},
	{600.0, 10.0},
}

func TestMinutes(t *testing.T) {
	for _, tt := range minutesTable {
		s := tt.n
		if s.Minutes() != tt.expected {
			logExpectedResult(s, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}

var nanosecondsTable = []struct {
	n        Second
	expected int64
}{
	{3.0, 3e+09},
	{6.0, 6e+09},
	{9.0, 9e+09},
	{12.0, 12e+09},
	{15.0, 15e+09},
}

func TestNanoseconds(t *testing.T) {
	for _, tt := range nanosecondsTable {
		s := tt.n
		if s.Nanoseconds() != tt.expected {
			logExpectedResult(s, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}

var secondsTable = []struct {
	n        Second
	expected float64
}{
	{1.0, 1.0},
	{2.0, 2.0},
	{150.0, 150.0},
	{360.0, 360.0},
	{9000.0, 9000.0},
}

func TestSeconds(t *testing.T) {
	for _, tt := range secondsTable {
		s := tt.n
		if s.Seconds() != tt.expected {
			logExpectedResult(s, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}

var stringTable = []struct {
	n        Second
	expected string
}{
	{1.1, "1.1s"},
	{300.0, "5m0s"},
	{3600.0, "1h0m0s"},
	{7200.0, "2h0m0s"},
	{86400.0, "24h0m0s"},
}

func TestString(t *testing.T) {
	for _, tt := range stringTable {
		s := tt.n
		if s.String() != tt.expected {
			logExpectedResult(s, tt.expected)
			t.Error("Result isn't accurate.")
		}
	}
}
