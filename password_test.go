package main

import "testing"

func TestHashPassword(t *testing.T) {
	password := *new(Password)
	expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	password.hashPassword("angryMonkey")
	hash := password.PasswordHash
	if hash != expectedHash {
		t.Errorf("Hash was incorrect, got: %s, want: %s.", hash, expectedHash)
	}
}

