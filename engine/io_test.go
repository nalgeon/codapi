package engine

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLimitedWriter(t *testing.T) {
	var b bytes.Buffer
	w := LimitWriter(&b, 5)

	{
		src := []byte{1, 2, 3}
		n, err := w.Write(src)
		if n != 3 {
			t.Fatalf("write(1,2,3): expected n = 3, got %d", n)
		}
		if err != nil {
			t.Fatalf("write(1,2,3): expected nil err, got %v", err)
		}
		if !reflect.DeepEqual(b.Bytes(), src) {
			t.Fatalf("write(1,2,3): expected %v, got %v", src, b.Bytes())
		}
	}

	{
		src := []byte{4, 5}
		n, err := w.Write(src)
		if n != 2 {
			t.Fatalf("+write(4,5): expected n = 2, got %d", n)
		}
		if err != nil {
			t.Fatalf("+write(4,5): expected nil err, got %v", err)
		}
		want := []byte{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(b.Bytes(), want) {
			t.Fatalf("+write(4,5): expected %v, got %v", want, b.Bytes())
		}
	}

	{
		src := []byte{6, 7, 8}
		n, err := w.Write(src)
		if n != 3 {
			t.Fatalf("+write(6,7,8): expected n = 3, got %d", n)
		}
		if err != nil {
			t.Fatalf("+write(6,7,8): expected nil err, got %v", err)
		}
		want := []byte{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(b.Bytes(), want) {
			t.Fatalf("+write(6,7,8): expected %v, got %v", want, b.Bytes())
		}
	}
}
