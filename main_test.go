package main

import (
	"path/filepath"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	allowedIDs = [2]int64{111, 222}
	if !isAllowed(111) {
		t.Error("user 111 should be allowed")
	}
	if !isAllowed(222) {
		t.Error("user 222 should be allowed")
	}
	if isAllowed(333) { // V1: unknown user → drop
		t.Error("user 333 should not be allowed")
	}
}

func TestOtherUser(t *testing.T) {
	allowedIDs = [2]int64{111, 222}
	if otherUser(111) != 222 { // V3: sender ≠ recipient
		t.Error("other of 111 should be 222")
	}
	if otherUser(222) != 111 {
		t.Error("other of 222 should be 111")
	}
}

func TestCounterPersistRoundtrip(t *testing.T) { // V4
	path := filepath.Join(t.TempDir(), "counter.json")
	if err := saveCounter(path, 42); err != nil {
		t.Fatal(err)
	}
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Errorf("want 42, got %d", got)
	}
}

func TestCounterDefaultsToZero(t *testing.T) { // V4: missing file → 0
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Errorf("want 0, got %d", got)
	}
}

func TestCounterNonNegative(t *testing.T) { // V5
	path := filepath.Join(t.TempDir(), "counter.json")
	count := 0
	for i := 0; i < 5; i++ {
		count++
		if err := saveCounter(path, count); err != nil {
			t.Fatal(err)
		}
	}
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got < 0 {
		t.Errorf("counter must be ≥ 0, got %d", got)
	}
}
