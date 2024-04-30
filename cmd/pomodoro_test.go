package main

import "testing"

func TestPomodoro(t *testing.T) {
	type TT struct {
		ts          int
		count, prog int
		status      string
	}

	tests := []TT{
		{0, 0, 0, "w"},
		{1, 0, 4, "w"},
		{26, 1, 20, "b"},
		{30, 1, 0, "w"},
		{35, 1, 20, "w"},
		{116, 4, 6, "lb"},
		{130, 4, 0, "w"},
		{135, 4, 20, "w"},
		{292, 9, 8, "w"},
		{380, 12, 33, "lb"},
	}

	for _, tt := range tests {
		th := float64(tt.ts) / 60

		count, prog, status := pomodoro(th)
		if count != tt.count || prog != tt.prog || status != tt.status {
			t.Errorf("pomodoro(%d) = %d, %d, %s; want %d, %d, %s", tt.ts, count, prog, status, tt.count, tt.prog, tt.status)
		}
	}
}
