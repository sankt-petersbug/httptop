package httptop

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/rivo/tview"
)

func TestFooterUpdateStatus(t *testing.T) {
	footer := &Footer{
		bucket: make([]int, 2),
	}

	footer.updateStatus(&Message{
		Records: []Record{
			Record{},
		},
	})

	if expected := 1; footer.total != expected {
		t.Errorf("expected total: %d, saw: %d", expected, footer.total)
	}
	if expected := []int{0, 1}; !reflect.DeepEqual(expected, footer.bucket) {
		t.Errorf("expected bucket: %v, saw: %v", expected, footer.bucket)
	}

	footer.updateStatus(&Message{
		Records: []Record{},
	})

	if expected := 1; footer.total != expected {
		t.Errorf("expected total: %d, saw: %d", expected, footer.total)
	}
	if expected := []int{1, 0}; !reflect.DeepEqual(expected, footer.bucket) {
		t.Errorf("expected bucket: %v, saw: %v", expected, footer.bucket)
	}
}

func TestFooterShouldAlert(t *testing.T) {
	testCases := []struct {
		name      string
		total     int
		rateLimit int
		expected  bool
	}{
		{
			name:      "No Alert",
			total:     0,
			rateLimit: 1,
			expected:  false,
		},
		{
			name:      "Alert",
			total:     121,
			rateLimit: 1,
			expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			footer := &Footer{
				total:     tc.total,
				rateLimit: tc.rateLimit,
			}
			result := footer.shouldAlert()

			if result != tc.expected {
				t.Errorf("expected result: %v, saw: %v", tc.expected, result)
			}
		})
	}
}

func TestFooterAlert(t *testing.T) {
	now := time.Now()
	footer := &Footer{
		view:      tview.NewTable().SetFixed(1, 2),
		bucket:    make([]int, 12),
		rateLimit: 0,
	}
	footer.Init()

	expected := "No Alert Detected"
	if cell := footer.view.GetCell(0, 1); !strings.Contains(cell.Text, expected) {
		t.Errorf("expected text: %q, saw: %q", expected, cell.Text)
	}

	footer.Update(&Message{
		Records: []Record{
			Record{},
		},
		CreatedAt: now,
	})

	expected = "High traffic generated an alert - hits = 1, triggered at"
	if cell := footer.view.GetCell(0, 1); !strings.Contains(cell.Text, expected) {
		t.Errorf("expected text: %q, saw: %q", expected, cell.Text)
	}
}
