package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type State struct {
	mu     sync.RWMutex
	states map[string]string
}

func New() *State {
	return &State{
		states: make(map[string]string),
	}
}

func Load(path string) (*State, error) {
	s := New()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return s, nil
	}

	if err := json.Unmarshal(data, &s.states); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *State) Get(feedID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.states[feedID]
}

func (s *State) Set(feedID, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[feedID] = value
}

func (s *State) Save(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.states, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}
