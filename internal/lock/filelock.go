package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileLock represents a file-based lock for preventing concurrent runs
type FileLock struct {
	lockPath string
	file     *os.File
}

// NewFileLock creates a new file lock instance
func NewFileLock(lockDir string) *FileLock {
	if lockDir == "" {
		// Default to temp directory if not specified
		lockDir = os.TempDir()
	}
	lockPath := filepath.Join(lockDir, "terradrift-watcher.lock")
	return &FileLock{
		lockPath: lockPath,
	}
}

// Acquire attempts to acquire the lock
func (fl *FileLock) Acquire() error {
	// Check if lock file exists and if it's stale
	if info, err := os.Stat(fl.lockPath); err == nil {
		// Lock file exists, check if it's stale (older than 1 hour)
		if time.Since(info.ModTime()) > time.Hour {
			// Stale lock, try to remove it
			os.Remove(fl.lockPath)
		} else {
			// Lock is fresh, another instance is running
			return fmt.Errorf("another instance is already running (lock file: %s)", fl.lockPath)
		}
	}

	// Try to create the lock file exclusively
	file, err := os.OpenFile(fl.lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("another instance is already running (lock file: %s)", fl.lockPath)
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	fl.file = file

	// Write PID and timestamp to the lock file
	lockInfo := fmt.Sprintf("PID: %d\nTime: %s\n", os.Getpid(), time.Now().Format(time.RFC3339))
	file.WriteString(lockInfo)

	return nil
}

// Release releases the lock
func (fl *FileLock) Release() error {
	if fl.file != nil {
		fl.file.Close()
		fl.file = nil
	}

	// Remove the lock file
	if err := os.Remove(fl.lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}

	return nil
}

// ForceRelease forcefully releases a lock (used with --force flag)
func (fl *FileLock) ForceRelease() error {
	// Close file if open
	if fl.file != nil {
		fl.file.Close()
		fl.file = nil
	}

	// Force remove the lock file
	if err := os.Remove(fl.lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to force remove lock file: %w", err)
	}

	return nil
}
