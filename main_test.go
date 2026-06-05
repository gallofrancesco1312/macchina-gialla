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
	if isAllowed(333) { // V1
		t.Error("user 333 should not be allowed")
	}
}

func TestCounterPersistRoundtrip(t *testing.T) { // V4
	path := filepath.Join(t.TempDir(), "counter.json")
	in := counterStore{CountA: 5, CountB: 3}
	if err := saveCounter(path, in); err != nil {
		t.Fatal(err)
	}
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != in {
		t.Errorf("want %+v, got %+v", in, got)
	}
}

func TestCounterDefaultsToZero(t *testing.T) { // V4: missing file → 0
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.CountA != 0 || got.CountB != 0 {
		t.Errorf("want {0,0}, got %+v", got)
	}
}

func TestCounterNonNegative(t *testing.T) { // V5
	path := filepath.Join(t.TempDir(), "counter.json")
	counts := counterStore{}
	for i := 0; i < 3; i++ {
		counts.CountA++
		counts.CountB++
	}
	if err := saveCounter(path, counts); err != nil {
		t.Fatal(err)
	}
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.CountA < 0 || got.CountB < 0 {
		t.Errorf("counters must be ≥ 0, got %+v", got)
	}
}

func TestPerUserCounters(t *testing.T) { // V6
	path := filepath.Join(t.TempDir(), "counter.json")
	counts := counterStore{CountA: 7, CountB: 2}
	if err := saveCounter(path, counts); err != nil {
		t.Fatal(err)
	}
	got, err := loadCounter(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.CountA != 7 {
		t.Errorf("CountA: want 7, got %d", got.CountA)
	}
	if got.CountB != 2 {
		t.Errorf("CountB: want 2, got %d", got.CountB)
	}
}

func TestFormatMsg(t *testing.T) {
	got := formatMsg("Macchina gialla! Totale: {count}", 42)
	want := "Macchina gialla! Totale: 42"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestFormatMsgNoPlaceholder(t *testing.T) {
	got := formatMsg("Bravissimo!", 10)
	if got != "Bravissimo!" {
		t.Errorf("unexpected change: %q", got)
	}
}
