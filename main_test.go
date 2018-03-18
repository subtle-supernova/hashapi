package main

import "testing"
import "time"
import "strings"
import "net/http"
import "net/http/httptest"

// TODO gofmt

func TestHashPassword(t *testing.T) {
    expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
    hash := hashPassword("angryMonkey")
    if hash != expectedHash {
       t.Errorf("Hash was incorrect, got: %s, want: %s.", hash, expectedHash)
    }
}

func TestSleepTimeSeconds(t *testing.T) {

    expectedVal := time.Duration(5)
    val := sleepTimeSeconds()
    if val != expectedVal {
       t.Errorf("Val was incorrect, got: %s, want: %s.", val, expectedVal)
    }
}

func TestNewShutdownValue(t *testing.T) {

    expectedVal := true
    val := newShutdownValue()
    if val != expectedVal {
       t.Errorf("Val was incorrect, got: %s, want: %s.", val, expectedVal)
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

func TestHash(t *testing.T) {
    resetVariablesToStartingValues()
    req, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
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

func TestHashAfterShutdown(t *testing.T) {
    resetVariablesToStartingValues()
    shutdownreq, err := http.NewRequest("GET", "/shutdown", nil)
    if err != nil {
        t.Fatal(err)
    }

    req, err := http.NewRequest("POST", "/hash", strings.NewReader("password=angryMonkey"))
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

