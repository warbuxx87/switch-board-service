package switch_board_service

import "testing"

func TestSample(t *testing.T) {
	total := 10
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}

