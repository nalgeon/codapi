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
func (m *Memory) Has(msg string) bool {
	for _, line := range m.Lines {
		if strings.Contains(line, msg) {
			return true
		}
	}
	return false
}

// MustHave checks if the memory has the message.
func (m *Memory) MustHave(t *testing.T, msg string) {
	if !m.Has(msg) {
		t.Errorf("%s must have: %s", m.Name, msg)
	}
}

// MustNotHave checks if the memory does not have the message.
func (m *Memory) MustNotHave(t *testing.T, msg string) {
	if m.Has(msg) {
		t.Errorf("%s must NOT have: %s", m.Name, msg)
	}
}

// Clear cleares the memory.
func (m *Memory) Clear() {
	m.Lines = []string{}
}

// Print prints memory lines to stdout.
func (m *Memory) Print() {
	for _, line := range m.Lines {
		fmt.Println(line)
	}
}
