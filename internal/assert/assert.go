package assert

import "testing"

// Equal helper compares two values of the same type that are comparable with == operator
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
