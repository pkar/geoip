package geoip

import "testing"

func Test_ipToInt(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int
		wantErr error
	}{
		{"", "174.36.207.186", 2921648058, nil},
		{"", "s.36.207.186", 0, ErrInvalidIP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipToInt(tt.in)
			equals(t, tt.wantErr, err)
			equals(t, tt.want, got)
		})
	}
}

func Test_intToIP(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		want    string
		wantErr error
	}{
		{"", 2921648058, "174.36.207.186", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := intToIP(tt.in)
			equals(t, tt.wantErr, err)
			equals(t, tt.want, got)
		})
	}
}
