package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasProperty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "present value",
			value: "Mallocs",
			want:  true,
		},
		{
			name:  "case sensitive",
			value: "mallocs",
			want:  false,
		},
		{
			name:  "spaces sensitive",
			value: "Mallocs ",
			want:  false,
		},
		{
			name:  "missing value",
			value: "some",
			want:  false,
		},
	}
	m := Metrics{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, m.HasProperty(test.value))
		})
	}
}

func TestHasType(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "present value",
			value: "gauge",
			want:  true,
		},
		{
			name:  "case sensitive",
			value: "Gauge",
			want:  false,
		},
		{
			name:  "spaces sensitive",
			value: "gauge ",
			want:  false,
		},
		{
			name:  "missing value",
			value: "some",
			want:  false,
		},
	}
	m := Metrics{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, m.HasType(test.value))
		})
	}
}
