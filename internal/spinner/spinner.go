package spinner

import (
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var frames = []string{".", "..", "...", "....", ".....", "......"}

// Spinner provides visual feedback for long-running operations.
type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

// New creates a new spinner with the given message.
// Only displays if stderr is a TTY.
func New(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation.
func (s *Spinner) Start() {
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		// Non-TTY: just print the message once.
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", s.message)
		close(s.done)
		return
	}

	go func() {
		defer close(s.done)
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the line.
				_, _ = fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			default:
				_, _ = fmt.Fprintf(os.Stderr, "\r  %s %s", s.message, frames[i%len(frames)])
				i++
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner and optionally prints a final message.
func (s *Spinner) Stop(finalMessage string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.stop:
		// Already stopped.
	default:
		close(s.stop)
	}
	<-s.done

	if finalMessage != "" {
		_, _ = fmt.Fprintf(os.Stderr, "  %s\n", finalMessage)
	}
}

// Update changes the spinner message while it's running.
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}
