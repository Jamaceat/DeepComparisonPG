package progress

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ProgressBar represents a progress bar with timing and status information
type ProgressBar struct {
	total       int64
	current     int64
	width       int
	description string
	startTime   time.Time
	lastUpdate  time.Time
	mutex       sync.RWMutex
	completed   bool
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64, description string) *ProgressBar {
	return &ProgressBar{
		total:       total,
		current:     0,
		width:       50,
		description: description,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
		completed:   false,
	}
}

// Update increments the progress bar by the specified amount
func (pb *ProgressBar) Update(increment int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current += increment
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.lastUpdate = time.Now()
	pb.render()
}

// SetProgress sets the current progress to a specific value
func (pb *ProgressBar) SetProgress(current int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current = current
	if pb.current > pb.total {
		pb.current = pb.total
	}
	pb.lastUpdate = time.Now()
	pb.render()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current = pb.total
	pb.completed = true
	pb.lastUpdate = time.Now()
	pb.render()
	fmt.Println() // New line after completion
}

// FinishWithMessage completes the progress bar with a custom message
func (pb *ProgressBar) FinishWithMessage(message string) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.current = pb.total
	pb.completed = true
	pb.lastUpdate = time.Now()

	elapsed := pb.lastUpdate.Sub(pb.startTime)
	// Clear line and show final completion message
	fmt.Printf("\r\033[K%s ✅ %s (completed in %v)\n",
		pb.description, message, formatDuration(elapsed))
}

// render draws the progress bar
func (pb *ProgressBar) render() {
	if pb.total == 0 {
		return
	}

	percentage := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * float64(pb.current) / float64(pb.total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)
	elapsed := pb.lastUpdate.Sub(pb.startTime)

	var eta string
	if pb.current > 0 && pb.current < pb.total {
		rate := float64(pb.current) / elapsed.Seconds()
		remaining := float64(pb.total-pb.current) / rate
		eta = fmt.Sprintf(" ETA: %v", formatDuration(time.Duration(remaining)*time.Second))
	} else if pb.completed {
		eta = fmt.Sprintf(" ✅ Completed in %v", formatDuration(elapsed))
	} else {
		eta = ""
	}

	// Clear line and show progress - use \r to return to beginning without newline
	status := fmt.Sprintf("\r\033[K%s [%s] %.1f%% (%d/%d) %v%s",
		pb.description, bar, percentage, pb.current, pb.total, formatDuration(elapsed), eta)

	fmt.Print(status)
	// Ensure we flush the output buffer
	os.Stdout.Sync()
}

// ConnectionProgress shows progress for database connection attempts
type ConnectionProgress struct {
	startTime   time.Time
	description string
	ticker      *time.Ticker
	done        chan bool
	mutex       sync.RWMutex
}

// NewConnectionProgress creates a new connection progress indicator
func NewConnectionProgress(description string) *ConnectionProgress {
	cp := &ConnectionProgress{
		startTime:   time.Now(),
		description: description,
		ticker:      time.NewTicker(100 * time.Millisecond),
		done:        make(chan bool),
	}

	go cp.animate()
	return cp
}

// animate shows a spinning animation during connection
func (cp *ConnectionProgress) animate() {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-cp.done:
			return
		case <-cp.ticker.C:
			cp.mutex.RLock()
			elapsed := time.Since(cp.startTime)
			// Clear line and show spinner animation
			fmt.Printf("\r\033[K%s %s %s", cp.description, spinner[i%len(spinner)], formatDuration(elapsed))
			os.Stdout.Sync()
			cp.mutex.RUnlock()
			i++
		}
	}
}

// Success completes the connection progress with success message
func (cp *ConnectionProgress) Success(message string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.ticker.Stop()
	close(cp.done)
	elapsed := time.Since(cp.startTime)
	// Clear line and show final success message
	fmt.Printf("\r\033[K%s ✅ %s (%v)\n", cp.description, message, formatDuration(elapsed))
	os.Stdout.Sync()
}

// Error completes the connection progress with error message
func (cp *ConnectionProgress) Error(message string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.ticker.Stop()
	close(cp.done)
	elapsed := time.Since(cp.startTime)
	// Clear line and show final error message
	fmt.Printf("\r\033[K%s ❌ %s (%v)\n", cp.description, message, formatDuration(elapsed))
	os.Stdout.Sync()
}

// SimpleProgress shows simple progress updates for operations without known total
type SimpleProgress struct {
	startTime   time.Time
	description string
	current     int64
	mutex       sync.RWMutex
}

// NewSimpleProgress creates a new simple progress indicator
func NewSimpleProgress(description string) *SimpleProgress {
	return &SimpleProgress{
		startTime:   time.Now(),
		description: description,
		current:     0,
	}
}

// Update increments the simple progress counter
func (sp *SimpleProgress) Update(increment int64) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	sp.current += increment
	elapsed := time.Since(sp.startTime)

	// Clear line and show simple progress
	fmt.Printf("\r\033[K%s: %d processed (%v)",
		sp.description, sp.current, formatDuration(elapsed))
	os.Stdout.Sync()
}

// Finish completes the simple progress with final count
func (sp *SimpleProgress) Finish(message string) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	elapsed := time.Since(sp.startTime)
	// Clear line and show final completion message
	fmt.Printf("\r\033[K%s ✅ %s - %d processed (%v)\n",
		sp.description, message, sp.current, formatDuration(elapsed))
	os.Stdout.Sync()
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// ShowProgress is a utility function for simple progress display
func ShowProgress(current, total int64, description string) {
	if total == 0 {
		return
	}

	percentage := float64(current) / float64(total) * 100
	filled := int(percentage / 2) // 50 character width

	bar := strings.Repeat("█", filled) + strings.Repeat("░", 50-filled)

	// Clear line and show progress
	fmt.Printf("\r\033[K%s [%s] %.1f%% (%d/%d)",
		description, bar, percentage, current, total)
	os.Stdout.Sync()

	if current >= total {
		fmt.Println(" ✅")
		os.Stdout.Sync()
	}
}
