package logx

import (
	"fmt"
	"strings"
	"testing"
)

// Memory stores logged messages in a slice.
type Memory struct {
	Name  string
	Lines []string
}

// NewMemory creates a new memory destination.
func NewMemory(name string) *Memory {
	return &Memory{Name: name, Lines: []string{}}
}

// Write implements the io.Writer interface.
func (m *Memory) Write(p []byte) (n int, err error) {
	msg := string(p)
	m.Lines = append(m.Lines, msg)
	return len(p), nil
}

// WriteString writes a string to the memory.
func (m *Memory) WriteString(s string) {
	m.Lines = append(m.Lines, s)
}

// Has returns true if the memory has the message.
func (m *Memory) Has(message ...string) bool {
	for _, line := range m.Lines {
		containsAll := true
		for _, part := range message {
			if !strings.Contains(line, part) {
				containsAll = false
				break
			}
		}
		if containsAll {
			return true
		}
	}
	return false
}

// MustHave checks if the memory has the message.
// If the message consists of several parts,
// they must all be in the same memory line.
func (m *Memory) MustHave(t *testing.T, message ...string) {
	if !m.Has(message...) {
		t.Errorf("%s must have: %v", m.Name, message)
	}
}

// MustNotHave checks if the memory does not have the message.
func (m *Memory) MustNotHave(t *testing.T, message ...string) {
	if m.Has(message...) {
		t.Errorf("%s must NOT have: %v", m.Name, message)
	}
}

// Clear clears the memory.
func (m *Memory) Clear() {
	m.Lines = []string{}
}

// Print prints memory lines to stdout.
func (m *Memory) Print() {
	for _, line := range m.Lines {
		fmt.Println(line)
	}
}
