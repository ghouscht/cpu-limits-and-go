package main

import "testing"

func TestIsPrime(t *testing.T) {
	tests := []struct {
		input uint64
		want  bool
	}{
		{
			input: 1,
			want:  false,
		},
		{
			input: 2,
			want:  true,
		},
		{
			input: 3,
			want:  true,
		},
		{
			input: 4,
			want:  false,
		},
		{
			input: 5,
			want:  true,
		},
	}

	for _, tc := range tests {
		got := isPrime(tc.input)
		if got != tc.want {
			t.Fatalf("isPrime returned %v but wanted %v for %d", got, tc.want, tc.input)
		}
	}
}
