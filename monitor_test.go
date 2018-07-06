package httptop

import (
	"reflect"
	"testing"
)

func TestGetSection(t *testing.T) {
	testCases := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "single path",
			request:  "/report",
			expected: "/report",
		},
		{
			name:     "double paths",
			request:  "/api/users",
			expected: "/api",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getSection(tc.request)

			if result != tc.expected {
				t.Errorf("expected result: %q, saw: %q", tc.expected, result)
			}
		})
	}
}

func TestGetTopHits(t *testing.T) {
	testCases := []struct {
		name     string
		records  []Record
		limit    int
		expected []SectionHits
	}{
		{
			name: "count sections",
			records: []Record{
				Record{Request: "/api/users"},
				Record{Request: "/api/instances"},
			},
			limit: 1,
			expected: []SectionHits{
				SectionHits{
					Section: "/api",
					Hits:    2,
				},
			},
		},
		{
			name: "sort by hits DESC",
			records: []Record{
				Record{Request: "/api/users"},
				Record{Request: "/reports"},
				Record{Request: "/api/instances"},
			},
			limit: 2,
			expected: []SectionHits{
				SectionHits{
					Section: "/api",
					Hits:    2,
				},
				SectionHits{
					Section: "/reports",
					Hits:    1,
				},
			},
		},
		{
			name: "limits",
			records: []Record{
				Record{Request: "/api/users"},
				Record{Request: "/reports"},
				Record{Request: "/api/instances"},
			},
			limit: 1,
			expected: []SectionHits{
				SectionHits{
					Section: "/api",
					Hits:    2,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetTopHits(tc.records, tc.limit)

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected result: %#v, saw: %#v", tc.expected, result)
			}
		})
	}
}
