package main

import "testing"
import "time"


func TestHashPassword(t *testing.T) {
    expected_hash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
    hash := hashPassword("angryMonkey")
    if hash != expected_hash {
       t.Errorf("Hash was incorrect, got: %s, want: %s.", hash, expected_hash)
    }
}

func TestSleepTimeSeconds(t *testing.T) {

    expected_val := time.Duration(5)
    val := sleepTimeSeconds()
    if val != expected_val {
       t.Errorf("Val was incorrect, got: %s, want: %s.", val, expected_val)
    }
}

func TestNewShutdownValue(t *testing.T) {

    expected_val := true
    val := newShutdownValue()
    if val != expected_val {
       t.Errorf("Val was incorrect, got: %s, want: %s.", val, expected_val)
    }
}

