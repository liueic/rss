package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestState_GetSet(t *testing.T) {
	s := New()

	if got := s.Get("test"); got != "" {
		t.Errorf("Get() = %v, want empty string", got)
	}

	s.Set("test", "value1")
	if got := s.Get("test"); got != "value1" {
		t.Errorf("Get() = %v, want value1", got)
	}

	s.Set("test", "value2")
	if got := s.Get("test"); got != "value2" {
		t.Errorf("Get() = %v, want value2", got)
	}
}

func TestState_SaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "test_state.json")

	s1 := New()
	s1.Set("feed1", "item1")
	s1.Set("feed2", "item2")

	if err := s1.Save(statePath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	s2, err := Load(statePath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := s2.Get("feed1"); got != "item1" {
		t.Errorf("Get(feed1) = %v, want item1", got)
	}

	if got := s2.Get("feed2"); got != "item2" {
		t.Errorf("Get(feed2) = %v, want item2", got)
	}
}

func TestState_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "nonexistent.json")

	s, err := Load(statePath)
	if err != nil {
		t.Fatalf("Load() error = %v, expected nil", err)
	}

	if s == nil {
		t.Fatal("Load() returned nil state")
	}

	if got := s.Get("test"); got != "" {
		t.Errorf("Get() = %v, want empty string", got)
	}
}

func TestState_LoadEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "empty.json")

	if err := os.WriteFile(statePath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	s, err := Load(statePath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if s == nil {
		t.Fatal("Load() returned nil state")
	}
}
