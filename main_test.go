package main

import "testing"


func TestHash(t *testing.T) {
    expected_hash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
    hash := hash("angryMonkey")
    if hash != expected_hash {
       t.Errorf("Hash was incorrect, got: %s, want: %s.", hash, expected_hash)
    }
}

