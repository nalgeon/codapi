// Package logx provides helper functions for logging.
package logx

import (
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)
var Verbose = false

// SetOutput sets the output destination.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// Printf prints a formatted message.
func Printf(format string, v ...any) {
	logger.Printf(format, v...)
}

// Println prints a message.
func Println(v ...any) {
	logger.Println(v...)
}

// Log prints a message.
func Log(message string, args ...any) {
	if len(args) == 0 {
		logger.Println(message)
	} else {
		logger.Printf(message+"\n", args...)
	}
}

// Debug prints a message if the verbose mode is on.
func Debug(message string, args ...any) {
	if !Verbose {
		return
	}
	Log(message, args...)
}

// Mock creates a new Memory and installs it as the logger output
// instead of the default one. Should be used for testing purposes only.
func Mock(path ...string) *Memory {
	memory := NewMemory("log")
	SetOutput(memory)
	Verbose = true
	return memory
}
