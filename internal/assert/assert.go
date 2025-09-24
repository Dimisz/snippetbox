package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	//indicates to Go test runner that it's a tester
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
