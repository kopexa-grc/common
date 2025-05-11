// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package graceful

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Ensure MockShutdownable implements Shutdownable interface
var _ Shutdownable = (*MockShutdownable)(nil)

// MockShutdownable ist eine Test-Implementierung von Shutdownable
type MockShutdownable struct {
	ShutdownCalled atomic.Bool
	ShutdownDelay  time.Duration
	ShouldError    bool
}

func (m *MockShutdownable) Shutdown(ctx context.Context) error {
	m.ShutdownCalled.Store(true)

	if m.ShutdownDelay > 0 {
		select {
		case <-time.After(m.ShutdownDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if m.ShouldError {
		return assert.AnError
	}

	return nil
}

func TestCloser_Register(t *testing.T) {
	closer := NewCloser()
	closer.SetSkipExit(true)

	mock := &MockShutdownable{}
	closer.Register("test", mock, 1*time.Second)

	targets := closer.GetTargets()
	assert.Len(t, targets, 1)
	assert.Equal(t, "test", targets[0].name)
	assert.Equal(t, time.Second, targets[0].timeout)
}

func TestCloser_DetectShutdown(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *MockShutdownable
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Normaler Shutdown",
			setupMock: func() *MockShutdownable {
				return &MockShutdownable{}
			},
			timeout: time.Second,
		},
		{
			name: "Shutdown mit Verzögerung",
			setupMock: func() *MockShutdownable {
				return &MockShutdownable{
					ShutdownDelay: 100 * time.Millisecond,
				}
			},
			timeout: time.Second,
		},
		{
			name: "Timeout beim Shutdown",
			setupMock: func() *MockShutdownable {
				return &MockShutdownable{
					ShutdownDelay: 2 * time.Second,
				}
			},
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closer := NewCloser()
			closer.SetSkipExit(true)

			mock := tt.setupMock()
			closer.Register("test", mock, tt.timeout)

			// Starte den Shutdown-Detector und warte auf Bereitschaft
			trigger, ready := closer.DetectShutdown()
			<-ready // Warte, bis der Detector bereit ist

			// Löse den Shutdown aus
			trigger()

			// Warte etwas länger als die Timeout-Dauer
			time.Sleep(tt.timeout + 100*time.Millisecond)

			assert.True(t, mock.ShutdownCalled.Load(), "Shutdown sollte aufgerufen worden sein")
		})
	}
}

func TestCloser_MultipleTargets(t *testing.T) {
	closer := NewCloser()
	closer.SetSkipExit(true)

	// Erstelle mehrere Mocks mit unterschiedlichen Verzögerungen
	mock1 := &MockShutdownable{ShutdownDelay: 50 * time.Millisecond}
	mock2 := &MockShutdownable{ShutdownDelay: 100 * time.Millisecond}
	mock3 := &MockShutdownable{ShutdownDelay: 150 * time.Millisecond}

	closer.Register("test1", mock1, time.Second)
	closer.Register("test2", mock2, time.Second)
	closer.Register("test3", mock3, time.Second)

	trigger, ready := closer.DetectShutdown()
	<-ready // Warte, bis der Detector bereit ist

	// Löse den Shutdown aus
	trigger()

	// Warte lange genug, damit alle Shutdowns abgeschlossen sein sollten
	time.Sleep(300 * time.Millisecond)

	assert.True(t, mock1.ShutdownCalled.Load(), "Shutdown 1 sollte aufgerufen worden sein")
	assert.True(t, mock2.ShutdownCalled.Load(), "Shutdown 2 sollte aufgerufen worden sein")
	assert.True(t, mock3.ShutdownCalled.Load(), "Shutdown 3 sollte aufgerufen worden sein")
}

func TestCloser_ConcurrentRegistration(t *testing.T) {
	closer := NewCloser()
	closer.SetSkipExit(true)

	// Registriere mehrere Targets gleichzeitig
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			mock := &MockShutdownable{}
			closer.Register(fmt.Sprintf("test%d", id), mock, time.Second)
		}(i)
	}

	wg.Wait()

	targets := closer.GetTargets()
	assert.Len(t, targets, 10)
}

func TestCloser_ShutdownError(t *testing.T) {
	closer := NewCloser()
	closer.SetSkipExit(true)

	mock := &MockShutdownable{ShouldError: true}
	closer.Register("test", mock, time.Second)

	trigger, ready := closer.DetectShutdown()
	<-ready // Warte, bis der Detector bereit ist

	trigger()
	time.Sleep(100 * time.Millisecond)

	assert.True(t, mock.ShutdownCalled.Load(), "Shutdown sollte trotz Fehler aufgerufen worden sein")
}

func TestCloser_ShutdownTimeout(t *testing.T) {
	closer := NewCloser()
	closer.SetSkipExit(true)

	mock := &MockShutdownable{ShutdownDelay: 200 * time.Millisecond}
	closer.Register("test", mock, 100*time.Millisecond) // Timeout kürzer als Delay

	trigger, ready := closer.DetectShutdown()
	<-ready

	trigger()
	time.Sleep(300 * time.Millisecond) // Warte länger als der Delay

	assert.True(t, mock.ShutdownCalled.Load(), "Shutdown sollte trotz Timeout aufgerufen worden sein")
}
