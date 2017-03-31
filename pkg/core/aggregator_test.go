package core

import (
	"testing"

	"github.com/taku-k/u2s3/pkg"
)

func TestParseEpoch(t *testing.T) {
	etests := []struct {
		desc string
		in   string
		out  string
	}{
		{"base", "time:24/Feb/2017:10:00:07 +0900\thost:127.0.0.1", "20170224100000"},
		{"just before begin carried", "time:24/Feb/2017:10:29:59 +0900\thost:127.0.0.1", "20170224100000"},
		{"just after begin carried", "time:24/Feb/2017:10:30:00 +0900\thost:127.0.0.1", "20170224103000"},
		{"time is after host", "host:127.0.0.1\ttime:24/Feb/2017:10:00:00 +0900", "20170224100000"},
		{"time is closed in the bracket", "host:127.0.0.1\ttime:[24/Feb/2017:10:00:00 +0900]", "20170224100000"},
	}
	cfg := &pkg.UploadConfig{LogFormat: "tsv", Step: 30}
	a := NewAggregator(nil, cfg)
	for _, tt := range etests {
		s := a.parseEpoch(tt.in)
		if s != tt.out {
			t.Errorf("%q: parseEpoch(%q) => %s, wants %q", tt.desc, tt.in, s, tt.out)
		}
	}
}
