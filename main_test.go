package main

import "testing"
import "time"
import "strings"
import "net/http"
import "net/http/httptest"
import "encoding/json"

// import "log"

func TestStatsOutputEmpty(t *testing.T) {
	stat := *new(Statistics)

	expected := `{"average":0,"total":0}`
	if stat.statsOutput() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			stat.statsOutput(), expected)
	}

}

func TestSleepTimeSeconds(t *testing.T) {

	expectedVal := time.Duration(5)
	val := sleepTimeSeconds()
	if val != expectedVal {
		t.Errorf("Val was incorrect, got: %s, want: %s.", val, expectedVal)
	}
}

func TestGetIdFromPathWithVal(t *testing.T) {

	expectedVal := 100
	val, err := getIdPointerFromPath("/hash/100")
	if err != nil {
		t.Errorf("Val was incorrect, got: nil, want: %v.", expectedVal)
		return
	}

	dereferencedVal := *val
	if dereferencedVal != expectedVal {
		t.Errorf("Val was incorrect, got: %v, want: %v.", dereferencedVal, expectedVal)
	}
}

func TestGetIdFromPathWithNoVal(t *testing.T) {

	expectedVal := (*int)(nil)
	val, _ := getIdPointerFromPath("/hash")
	if val != expectedVal {
		t.Errorf("Val was incorrect, got: %v, want: %v.", val, expectedVal)
	}
}

func TestRegisterShutdown(t *testing.T) {
	resetVariablesToStartingValues()
	req, err := http.NewRequest("GET", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerShutdown)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusOK
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}
}

func TestHashA(t *testing.T) {
	resetVariablesToStartingValues()
	req, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusOK
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}
}

func TestHashNoPassword(t *testing.T) {
	resetVariablesToStartingValues()
	req, err := http.NewRequest("POST", "/hash", strings.NewReader("pasword=angryMonkey"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusBadRequest
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

}

func TestHashAfterShutdown(t *testing.T) {
	resetVariablesToStartingValues()
	shutdownreq, err := http.NewRequest("GET", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	shutdownhandler := http.HandlerFunc(registerShutdown)

	shutdownhandler.ServeHTTP(rr, shutdownreq)

	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusServiceUnavailable
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

}

func TestHashWithRequestIdNotAlreadySet(t *testing.T) {
	resetVariablesToStartingValues()
	req, err := http.NewRequest("GET", "/hash/100", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusNotFound
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

}

func TestHashWithRequestIdAlreadySetButPasswordNotHashed(t *testing.T) {
	resetVariablesToStartingValues()
	initialreq, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
	initialreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/hash/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	initialrr := httptest.NewRecorder()
	initialhandler := http.HandlerFunc(hash)

	initialhandler.ServeHTTP(initialrr, initialreq)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusNotFound
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

}

func TestHashWithRequestIdAlreadySetPasswordHashed(t *testing.T) {
	resetVariablesToStartingValues()
	initialreq, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
	initialreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/hash/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	initialrr := httptest.NewRecorder()
	initialhandler := http.HandlerFunc(hash)

	initialhandler.ServeHTTP(initialrr, initialreq)

	// wait for password to be created
	time.Sleep((1 + sleepTimeSeconds()) * time.Second)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hash)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusOK
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

}

func TestStatisticsEmpty(t *testing.T) {
	resetVariablesToStartingValues()
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statisticsGet)

	handler.ServeHTTP(rr, req)

	expectedVal := http.StatusOK
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

	expected := `{"average":0,"total":0}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestStatisticsOneRequest(t *testing.T) {
	resetVariablesToStartingValues()
	initialreq, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
	initialreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	initialhandler := http.HandlerFunc(hash)

	initialhandler.ServeHTTP(rr, initialreq)

	// wait for password to be created
	time.Sleep((1 + sleepTimeSeconds()) * time.Second)

	rrForStats := httptest.NewRecorder()
	handler := http.HandlerFunc(statisticsGet)

	handler.ServeHTTP(rrForStats, req)

	expectedVal := http.StatusOK
	if status := rr.Code; status != expectedVal {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedVal)
	}

	var dat map[string]int
	bodyStr := rrForStats.Body.String()
	body := []byte(bodyStr)

	if err := json.Unmarshal(body, &dat); err != nil {
		t.Errorf("handler returned unexpected json: %v", bodyStr)
	}

	expectedTotal := 1
	unexpectedAverage := 0
	if dat["total"] != expectedTotal {
		t.Errorf("handler returned unexpected total: got %v want %v",
			dat["total"], expectedTotal)
	}
	if dat["average"] == unexpectedAverage {
		t.Errorf("handler returned unexpected avreage: got %v don't want %v",
			dat["average"], unexpectedAverage)
	}

}
