package core

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	gzip "github.com/klauspost/pgzip"
)

func TestFileCompress(t *testing.T) {
	contents := []byte("abc\ndef\nghi")
	src, err := ioutil.TempFile("", "u2s3")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	src.Write(contents)
	src.Close()
	defer os.Remove(src.Name())
	f := NewFile(src.Name(), "", "", "")
	if err := f.compress(); err != nil {
		t.Fatalf("error: %s", err)
	}
	out := f.GetFile()
	defer os.Remove(out.Name())
	outRaw, err := gzip.NewReader(out)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	got, err := ioutil.ReadAll(outRaw)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !reflect.DeepEqual(contents, got) {
		t.Errorf("contents mismatch: %s, %s", contents, got)
	}
}
