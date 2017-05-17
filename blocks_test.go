package geoip

import (
	"strings"
	"testing"
)

func TestNewBlock(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    *Block
		wantErr error
	}{
		{
			"",
			[]string{"24100", "24999", "123"},
			&Block{StartIPNum: 24100, EndIPNum: 24999, LocationID: 123},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBlock(tt.in)
			equals(t, tt.wantErr, err)
			equals(t, tt.want, got)
		})
	}
}

func TestNewBlocksMap(t *testing.T) {
	r := strings.NewReader(`garbage
"16777216","16777471","609013"
"16777472","16778239","104084"
`)
	bm, err := NewBlocksMap(r)
	ok(t, err)
	got, ok := bm.Lookup("1.0.0.255")
	equals(t, 609013, got)
	equals(t, true, ok)

	got, ok = bm.Lookup("111")
	equals(t, 0, got)
	equals(t, false, ok)
}
