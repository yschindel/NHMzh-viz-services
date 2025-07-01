package azure

import (
	"testing"
	"time"
)

func TestTimestampRounding(t *testing.T) {
	cases := []struct {
		in   string
		want time.Time
	}{
		{
			in:   "2023-01-01T12:00:00.249Z",
			want: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC).Add(200 * time.Millisecond),
		},
		{
			in:   "2023-01-01T12:00:00.251Z",
			want: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC).Add(300 * time.Millisecond),
		},
	}

	for _, c := range cases {
		ts, err := time.Parse(time.RFC3339Nano, c.in)
		if err != nil {
			t.Fatalf("failed to parse timestamp %s: %v", c.in, err)
		}
		got := ts.Round(time.Millisecond * 100)
		if !got.Equal(c.want) {
			t.Errorf("rounded %s got %s want %s", c.in, got, c.want)
		}
	}
}
