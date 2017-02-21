package aggregator

import "testing"

func newAggregator() *Aggregator {
	return NewAggregator(nil, "tsv", "")
}

func TestParseEpoch(t *testing.T) {
	etests := []struct {
		desc string
		in   string
		out  string
	}{
		{"simple", "", ""},
	}

	a := newAggregator()
	for _, tt := range etests {
		s := a.parseEpoch(tt.in)
		if s != tt.out {
			t.Errorf("%q: parseEpoch(%q) => %s, wants %q", tt.desc, tt.in, s, tt.out)
		}
	}
}
