package assert

import (
	"fmt"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if want != got {
		panic(fmt.Sprintf("want %v; got %v", want, got))
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		panic(fmt.Sprintf("want %q to contain %q", actual, expectedSubstring))
	}
}
