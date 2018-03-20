package main

import "testing"

func TestStatsOutputOneValue(t *testing.T) {
	stat := *new(Statistics)
	stat.Total = 1
	stat.CumulativeTime = 100

	expected := `{"average":100,"total":1}`
	if stat.statsOutput() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			stat.statsOutput(), expected)
	}
}

func TestStatsOutputTwoValues(t *testing.T) {
	stat := *new(Statistics)
	stat.Total = 2
	stat.CumulativeTime = 300

	expected := `{"average":150,"total":2}`
	if stat.statsOutput() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			stat.statsOutput(), expected)
	}
}

func TestStatsIncrementTotal(t *testing.T) {
	stat := *new(Statistics)
	stat.incrementTotal()
	stat.incrementTotal()

	expected := 2
	if stat.Total != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			stat.statsOutput(), expected)
	}
}

func TestStatsIncrementCumulativeTime(t *testing.T) {
	stat := *new(Statistics)
	stat.incrementCumulativeTime(100)
	stat.incrementCumulativeTime(200)

	expected := 300
	if stat.CumulativeTime != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			stat.statsOutput(), expected)
	}
}
