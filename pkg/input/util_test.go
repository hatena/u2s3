package input

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsGzipped(t *testing.T) {
	cases := []struct {
		content string
		gzipped bool
	}{
		{"not compressed", false},
		{"compressed", true},
	}

	var w io.Writer
	for _, c := range cases {
		fp, err := ioutil.TempFile("", "u2s3-test")
		if err != nil {
			break
		}
		buf := bufio.NewWriter(fp)
		if c.gzipped {
			w = gzip.NewWriter(buf)
		} else {
			w = buf
		}
		w.Write([]byte(c.content))
		buf.Flush()
		fp.Seek(0, 0)
		if IsGzipped(fp) != c.gzipped {
			t.Errorf("IsGzipped => error: %q", c.gzipped)
		}
		fp.Close()
		os.Remove(fp.Name())
	}
}
