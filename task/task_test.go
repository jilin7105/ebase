package task

import (
	"testing"
	"time"
)

func TestInitTaskServer(t *testing.T) {
	s := InitTaskServer()
	if s == nil {
		t.Error("Expected task server, got nil")
	}

	_, err := s.Every(1).Second().Do(func() {
		t.Log("Task executed")
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	s.StartAsync()

	// Sleep for a while to allow the task to execute
	time.Sleep(2 * time.Second)

	s.Clear()
}
