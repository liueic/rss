package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	configYAML := `feeds:
  - id: test-feed
    name: Test Feed
    url: https://example.com/rss
    notify: true
    dedupe_key: guid
    aggregate: false
  - id: test-feed-2
    name: Test Feed 2
    url: https://example.com/rss2
    notify: false
    dedupe_key: link
    aggregate: true
    aggregate_window_minutes: 60
`

	if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Feeds) != 2 {
		t.Fatalf("Expected 2 feeds, got %d", len(cfg.Feeds))
	}

	feed1 := cfg.Feeds[0]
	if feed1.ID != "test-feed" {
		t.Errorf("feed1.ID = %v, want test-feed", feed1.ID)
	}
	if feed1.Name != "Test Feed" {
		t.Errorf("feed1.Name = %v, want Test Feed", feed1.Name)
	}
	if feed1.URL != "https://example.com/rss" {
		t.Errorf("feed1.URL = %v, want https://example.com/rss", feed1.URL)
	}
	if !feed1.Notify {
		t.Error("feed1.Notify = false, want true")
	}
	if feed1.DedupeKey != "guid" {
		t.Errorf("feed1.DedupeKey = %v, want guid", feed1.DedupeKey)
	}

	feed2 := cfg.Feeds[1]
	if feed2.Notify {
		t.Error("feed2.Notify = true, want false")
	}
	if feed2.Aggregate != true {
		t.Error("feed2.Aggregate = false, want true")
	}
	if feed2.AggregateWindowMinutes != 60 {
		t.Errorf("feed2.AggregateWindowMinutes = %v, want 60", feed2.AggregateWindowMinutes)
	}
}

func TestConfig_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() expected error, got nil")
	}
}
