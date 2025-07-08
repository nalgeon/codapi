package engine

import (
	"bytes"
	"testing"

	"github.com/nalgeon/be"
)

func TestLimitedWriter(t *testing.T) {
	var b bytes.Buffer
	w := LimitWriter(&b, 5)

	{
		src := []byte{1, 2, 3}
		n, err := w.Write(src)
		be.Err(t, err, nil)
		be.Equal(t, n, 3)
		be.Equal(t, b.Bytes(), src)
	}

	{
		src := []byte{4, 5}
		n, err := w.Write(src)
		be.Err(t, err, nil)
		be.Equal(t, n, 2)
		want := []byte{1, 2, 3, 4, 5}
		be.Equal(t, b.Bytes(), want)
	}

	{
		src := []byte{6, 7, 8}
		n, err := w.Write(src)
		be.Err(t, err, nil)
		be.Equal(t, n, 3)
		want := []byte{1, 2, 3, 4, 5}
		be.Equal(t, b.Bytes(), want)
	}
}
