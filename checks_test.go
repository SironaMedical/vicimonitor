package main

import (
	"testing"
)

func TestChildSAExists(t *testing.T) {
	cases := []struct {
		name        string
		sourceName  string
		sourceIkeSA *IkeSA
		expected    bool
	}{
		{
			name:        "empty",
			expected:    false,
			sourceName:  "",
			sourceIkeSA: &IkeSA{},
		},
		{
			name:        "0_ikeChild",
			expected:    false,
			sourceName:  "peer-2.2.2.2-tunnel-801",
			sourceIkeSA: &IkeSA{},
		},
		{
			name:       "match_1_ikeChild",
			expected:   true,
			sourceName: "peer-2.2.2.2-tunnel-801",
			sourceIkeSA: &IkeSA{
				ChildSAs: map[string]ChildSA{
					"peer-2.2.2.2-tunnel-801-1": {
						Name: "peer-2.2.2.2-tunnel-801",
					},
				},
			},
		},
		{
			name:       "no_match_1_ikeChild",
			expected:   false,
			sourceName: "peer-2.2.2.2-tunnel-802",
			sourceIkeSA: &IkeSA{
				ChildSAs: map[string]ChildSA{
					"peer-2.2.2.2-tunnel-801-1": {
						Name: "peer-2.2.2.2-tunnel-801",
					},
				},
			},
		},
		{
			name:       "match_2_ikeChild",
			expected:   true,
			sourceName: "peer-2.2.2.2-tunnel-801",
			sourceIkeSA: &IkeSA{
				ChildSAs: map[string]ChildSA{
					"peer-2.2.2.2-tunnel-801-1": {
						Name: "peer-2.2.2.2-tunnel-801",
					},
					"peer-2.2.2.2-tunnel-802-1": {
						Name: "peer-2.2.2.2-tunnel-802",
					},
				},
			},
		},
		{
			name:       "no_match_2_ikeChild",
			expected:   false,
			sourceName: "peer-2.2.2.2-tunnel-803",
			sourceIkeSA: &IkeSA{
				ChildSAs: map[string]ChildSA{
					"peer-2.2.2.2-tunnel-801-1": {
						Name: "peer-2.2.2.2-tunnel-801",
					},
					"peer-2.2.2.2-tunnel-802-1": {
						Name: "peer-2.2.2.2-tunnel-802",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ChildSAExists(tc.sourceName, tc.sourceIkeSA)

			if result != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}

func TestIkeSAEstablished(t *testing.T) {
	cases := []struct {
		name        string
		sourceIkeSA IkeSA
		expected    bool
	}{
		{
			name:        "empty",
			expected:    false,
			sourceIkeSA: IkeSA{},
		},
		{
			name:     "connecting_ikeSA",
			expected: false,
			sourceIkeSA: IkeSA{
				Name:  "peer-50.112.52.187-tunnel-801",
				State: "CONNECTING",
			},
		},
		{
			name:     "established_ikeSA",
			expected: true,
			sourceIkeSA: IkeSA{
				Name:  "peer-50.112.52.187-tunnel-801",
				State: "ESTABLISHED",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsIkeSAEstablished(tc.sourceIkeSA)

			if result != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
