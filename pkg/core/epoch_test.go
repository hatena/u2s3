package aggregator

import (
	"os"
	"testing"
)

func TestGetObjectKey(t *testing.T) {
	host, _ := os.Hostname()
	cases := []struct {
		desc   string
		key    string
		format string
		seq    int
		out    string
	}{
		{"base", "20170224173000", "{{.Output}}/{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.Minute}}_{{.Seq}}.log.gz", 1, "output/2017/02/24/" + host + "-201702241730_1.log.gz"},
		{"padding day", "20170204173000", "{{.Day}}", 1, "04"},
	}

	e, _ := NewEpoch("", "", "output")
	for _, c := range cases {
		e.epochKey = c.key
		e.keyFmt = c.format
		s, err := e.GetObjectKey(c.seq)
		if err != nil {
			t.Errorf("%q: GetObjectKey(%d) => error: %q", c.desc, c.seq, err)

		} else if s != c.out {
			t.Errorf("%q: GetObjectKey(%d) => %s, wants %q", c.desc, c.seq, s, c.out)
		}
	}
}
