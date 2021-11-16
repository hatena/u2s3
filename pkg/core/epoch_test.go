package core

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
		s, err := e.GetObjectKey()
		if err != nil {
			t.Errorf("%q: GetObjectKey() => error: %q", c.desc, err)

		} else if s != c.out {
			t.Errorf("%q: GetObjectKey() => %s, wants %q", c.desc, s, c.out)
		}
	}
}

func TestParseEpoch(t *testing.T) {
	etests := []struct {
		desc      string
		logFormat string
		in        string
		out       string
	}{
		{"base", "tsv", "time:24/Feb/2017:10:00:07 +0900\thost:127.0.0.1", "20170224100000"},
		{"just before begin carried", "tsv", "time:24/Feb/2017:10:29:59 +0900\thost:127.0.0.1", "20170224100000"},
		{"just after begin carried", "tsv", "time:24/Feb/2017:10:30:00 +0900\thost:127.0.0.1", "20170224103000"},
		{"time is after host", "tsv", "host:127.0.0.1\ttime:24/Feb/2017:10:00:00 +0900", "20170224100000"},
		{"time is closed in the bracket", "tsv", "host:127.0.0.1\ttime:[24/Feb/2017:10:00:00 +0900]", "20170224100000"},

		{"base", "jsonl", `{"time": "24/Feb/2017:10:00:07 +0900", "host": "127.0.0.1"}`, "20170224100000"},
		{"just before begin carried", "jsonl", `{"time": "24/Feb/2017:10:29:59 +0900", "host": "127.0.0.1"}`, "20170224100000"},
		{"just after begin carried", "jsonl", `{"time": "24/Feb/2017:10:30:00 +0900", "host": "127.0.0.1"}`, "20170224103000"},
		{"time is after host", "jsonl", `{"host": "127.0.0.1", "time": "24/Feb/2017:10:00:00 +0900"}`, "20170224100000"},
		{"no spacing", "jsonl", `{"host":"127.0.0.1","time":"24/Feb/2017:10:00:00 +0900"}`, "20170224100000"},
		{"lots of spacing", "jsonl", `{"host"   :  "127.0.0.1",    "time"   :  "24/Feb/2017:10:00:00 +0900"  }`, "20170224100000"},
	}
	for _, tt := range etests {
		s := parseEpoch(tt.in, tt.logFormat, 30)
		if s != tt.out {
			t.Errorf("[%q] %q: parseEpoch(%q) => %s, wants %q", tt.logFormat, tt.desc, tt.in, s, tt.out)
		}
	}
}
