package engine

import "io"

// A LimitedWriter writes to w but limits the amount
// of data to only n bytes. After reaching the limit,
// silently discards the rest of the data without errors.
type LimitedWriter struct {
	w io.Writer
	n int64
}

// LimitWriter returns a writer that writes no more
// than n bytes and silently discards the rest.
func LimitWriter(w io.Writer, n int64) io.Writer {
	return &LimitedWriter{w, n}
}

// Write implements the io.Writer interface.
func (w *LimitedWriter) Write(p []byte) (int, error) {
	lenp := len(p)
	if w.n <= 0 {
		return lenp, nil
	}
	if int64(lenp) > w.n {
		p = p[:w.n]
	}
	n, err := w.w.Write(p)
	w.n -= int64(n)
	return lenp, err
}
