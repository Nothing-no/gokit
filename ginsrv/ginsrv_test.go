package ginsrv

import "testing"

func TestCheckAddr(t *testing.T) {
	ts := []struct {
		addr string
		err  error
	}{
		{"123.123.231.34", nil},
		{"10.23.42.a", newErr},
		{"234.2133.23.4", newErr},
		{"10.23,42.32", newErr},
	}
	for _, tc := range ts {
		err := checkAddr(tc.addr)
		if nil != err {
			if err != tc.err {
				t.Errorf("want: %v, get: %v\n", tc.err, err)
			}
		}
	}
}
