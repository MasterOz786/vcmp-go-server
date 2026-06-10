package safari

import "testing"

func TestNormalizeCommand(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"/help", "/help", true},
		{"help", "/help", true},
		{"pack 1", "/pack 1", true},
		{"/pack 2", "/pack 2", true},
		{"  /status  ", "/status", true},
		{"", "", false},
		{"   ", "", false},
	}
	for _, tc := range tests {
		got, ok := normalizeCommand(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Fatalf("normalizeCommand(%q) = (%q, %v), want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}
