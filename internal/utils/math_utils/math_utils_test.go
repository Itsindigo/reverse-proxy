package math_utils

import (
	"testing"
)

func TestMin(t *testing.T) {
	got := Min(3, 9)
	want := 3

	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}
