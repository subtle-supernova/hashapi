package main

import "sync"
import "encoding/json"

type Statistics struct {
	Total          int
	CumulativeTime int
	mux            sync.Mutex
}

func (s *Statistics) statsOutput() string {
	total := s.Total
	averageTime := 0
	if total != 0 {
		averageTime = s.CumulativeTime / s.Total
	}
	statsMap := map[string]int{"total": total, "average": averageTime}
	jsonMap, _ := json.Marshal(statsMap)
	return string(jsonMap)
}

func (s *Statistics) incrementTotal() {
	s.mux.Lock()
	s.Total += 1
	s.mux.Unlock()
}

func (s *Statistics) incrementCumulativeTime(timeInMilliseconds int) {
	s.mux.Lock()
	s.CumulativeTime += timeInMilliseconds
	s.mux.Unlock()
}
