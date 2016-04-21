package requests

import (
	"time"
)

const nano = 1e-09

// Duration is an interface implemented by Second and time.Duration.
type Duration interface {
	Hours() float64
	Minutes() float64
	Nanoseconds() int64
	Seconds() float64
	String() string
}

// Second represents amount of seconds as a floating point number.
type Second float64

// Hours returns the duration as a floating point number of hours.
func (s Second) Hours() (t float64) {
	n := float64(s) / nano
	t = time.Duration(n).Hours()
	return
}

// Minutes returns the duration as a floating point number of minutes.
func (s Second) Minutes() (t float64) {
	n := float64(s) / nano
	t = time.Duration(n).Minutes()
	return
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (s Second) Nanoseconds() (ns int64) {
	n := float64(s) / nano
	ns = time.Duration(n).Nanoseconds()
	return
}

// Seconds returns the duration as a floating point number of seconds.
func (s Second) Seconds() float64 {
	return float64(s)
}

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0, with no unit.
func (s Second) String() (ts string) {
	n := float64(s) / nano
	ts = time.Duration(n).String()
	return
}

// Ftos returns the floating point number of seconds to duration.
func Ftos(s float64) (t time.Duration) {
	n := s / nano
	t = time.Duration(n)
	return
}

// Stof returns the duration as a floating point number of seconds.
func Stof(t time.Duration) (s float64) {
	s = t.Seconds()
	return
}
