package geoip

import (
	"strings"
	"testing"
)

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    *Location
		wantErr error
	}{
		{
			"",
			[]string{"24100", "CA", "QC", "Laval", "h7w4s8", "45.6167", "-73.7500", "", ""},
			Laval,
			nil,
		},
		{
			"",
			[]string{"24107", "US", "NY", "Fonda", "12068", "42.9508", "-74.3937", "532", "518"},
			Fonda,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLocation(tt.in)
			equals(t, tt.wantErr, err)
			equals(t, tt.want, got)
		})
	}
}

func TestNewLocationMap(t *testing.T) {
	r := strings.NewReader(`garbage
1,"O1","","","",0.0000,0.0000,,
24107,"US","NY","Fonda","12068",42.9508,-74.3937,532,518
`)
	lm, err := NewLocationMap(r)
	ok(t, err)
	got, ok := lm.LookupByID(24107)
	equals(t, Fonda, got)
	equals(t, true, ok)

	got, ok = lm.LookupByID(111)
	equals(t, (*Location)(nil), got)
	equals(t, false, ok)
}
