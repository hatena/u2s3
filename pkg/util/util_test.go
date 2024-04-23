package util

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"reflect"
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
		fp, err := os.CreateTemp("", "u2s3-test")
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
			t.Errorf("IsGzipped => error: %t", c.gzipped)
		}
		fp.Close()
		os.Remove(fp.Name())
	}
}

func TestGetParams(t *testing.T) {
	cases := []struct {
		regEx    string
		url      string
		expected map[string]string
	}{
		{"(?P<Year>\\d{4})(?P<Month>\\d{2})(?P<Day>\\d{2}).tsv", "20170331.tsv", map[string]string{"Year": "2017", "Month": "03", "Day": "31"}},
		{"", "/var/lib/mysql/mysqld-slow.log.1", map[string]string{}},
	}

	for _, c := range cases {
		regEx := c.regEx
		url := c.url
		actual := GetParams(regEx, url)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetParams => error: regEx: %s, url: %s, actual: %v, expected: %v", regEx, url, actual, c.expected)
		}
	}
}
