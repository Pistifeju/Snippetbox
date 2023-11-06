package main

import (
	"github.com/Pistieju/snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{
			name: "UTC",
			t:    time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
			want: "01 Jan 2021 at 12:00",
		},
		{
			name: "Empty",
			t:    time.Time{},
			want: "",
		},
		{
			name: "CET",
			t:    time.Date(2021, 1, 1, 12, 0, 0, 0, time.FixedZone("CET", 60*60)),
			want: "01 Jan 2021 at 11:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.t)
			assert.Equal(t, tt.want, got)
		})
	}
}
