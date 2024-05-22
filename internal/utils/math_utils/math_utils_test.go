package math_utils

import (
	"testing"
)

func TestToFixed(t *testing.T) {
	got := ToFixed(16.666666, 2)
	want := 16.67

	if got != want {
		t.Errorf("got '%f' want '%f'", got, want)
	}
}
