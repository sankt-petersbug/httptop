package httptop

import (
	"reflect"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected Record
	}{
		{
			name: "base",
			line: `127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 1234`,
			expected: Record{
				RemoteAddr: "127.0.0.1",
				Method:     "GET",
				Userid:     "james",
				Time:       time.Date(2018, time.May, 9, 16, 0, 39, 0, time.UTC),
				Request:    "/report",
				StatusCode: 200,
				Bytes:      1234,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewRecord(tc.line)

			if err != nil {
				t.Errorf("expected nil, but got an error %q", err)
			}

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected result: %v, saw: %v", tc.expected, result)
			}
		})
	}
}
