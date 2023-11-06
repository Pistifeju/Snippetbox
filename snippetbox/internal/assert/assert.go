package assert

import (
	"fmt"
	"testing"
)

func Equal[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if want != got {
		panic(fmt.Sprintf("want %v; got %v", want, got))
	}
}
