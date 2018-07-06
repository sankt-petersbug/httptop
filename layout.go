package httptop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rivo/tview"
	//"github.com/gdamore/tcell"
)

const noAlertMsg = "            No Alert Detected"

type Footer struct {
	view      *tview.Table
	noAlert   string
	bucket    []int
	rateLimit int
	total     int
}

func (c *Footer) Init() {
	c.view.
		SetCell(0, 0, tview.NewTableCell(time.Now().Format(time.RFC1123))).
		SetCell(0, 1, tview.NewTableCell(noAlertMsg))
}
func (c *Footer) updateStatus(msg *Message) {
	count := len(msg.Records)
	c.bucket = c.bucket[1:]
	c.bucket = append(c.bucket, count)

	c.total = 0
	for _, n := range c.bucket {
		c.total += n
	}
}
func (c *Footer) shouldAlert() bool {
	interval := len(c.bucket) * 10
	rate := float64(c.total) / float64(interval)

	return rate > float64(c.rateLimit)
}
func (c *Footer) Update(msg *Message) {
	ts := msg.CreatedAt.Format(time.RFC1123)
	alert := noAlertMsg

	c.updateStatus(msg)

	if c.shouldAlert() {
		alert = fmt.Sprintf("            High traffic generated an alert - hits = %d, triggered at %s", c.total, ts)
	}

	c.view.
		SetCell(0, 0, tview.NewTableCell(ts)).
		SetCell(0, 1, tview.NewTableCell(alert))
}

type MostHits struct {
	view *tview.Table
}

func (c *MostHits) Init() {
	header := []string{"Section            ", "Hits"}

	for i, text := range header {
		c.view.SetCell(0, i, tview.NewTableCell(text))
	}
}
func (c *MostHits) Update(msg *Message) {
	for i := 1; i < 6; i++ {
		c.view.
			SetCell(i, 0, tview.NewTableCell("")).
			SetCell(i, 1, tview.NewTableCell(""))
	}

	for i, item := range GetTopHits(msg.Records, 5) {
		row := i + 1

		c.view.
			SetCell(row, 0, tview.NewTableCell(item.Section)).
			SetCell(row, 1, tview.NewTableCell(strconv.Itoa(item.Hits)))
	}
}

type RequestCounts struct {
	view *tview.Table
}

func (c *RequestCounts) Init() {
	header := []string{"Total Requests    ", "Valid Requests    ", "Failed Requests    "}

	for i, text := range header {
		c.view.SetCell(i, 0, tview.NewTableCell(text))
	}

	for i := 0; i < 3; i++ {
		c.view.SetCell(i, 1, tview.NewTableCell("0"))
	}
}
func (c *RequestCounts) Update(msg *Message) {
	counts := []int{len(msg.Records), 0, 0}

	for _, record := range msg.Records {
		if 200 <= record.StatusCode && 400 > record.StatusCode {
			counts[1]++
		} else {
			counts[2]++
		}
	}

	for i, count := range counts {
		c.view.SetCell(i, 1, tview.NewTableCell(strconv.Itoa(count)))
	}
}

type UniqueVisitors struct {
	view *tview.Table
}

func (c *UniqueVisitors) Init() {
	c.view.
		SetCell(0, 0, tview.NewTableCell("Unique Visitors    ")).
		SetCell(0, 1, tview.NewTableCell("0")).
		SetCell(1, 0, tview.NewTableCell("Bytes Sent    ")).
		SetCell(1, 1, tview.NewTableCell("0"))
}
func (c *UniqueVisitors) Update(msg *Message) {
	m := make(map[string]bool)
	for _, record := range msg.Records {
		m[record.RemoteAddr] = true
	}

	c.view.SetCell(0, 1, tview.NewTableCell(strconv.Itoa(len(m))))

	total := 0
	for _, record := range msg.Records {
		total += record.Bytes
	}

	c.view.SetCell(1, 1, tview.NewTableCell(strconv.Itoa(total)))
}

type Component interface {
	Init()
	Update(msg *Message)
}

type Layout struct {
	view       tview.Primitive
	components []Component
}

func (l *Layout) Update(msg *Message) {
	for _, c := range l.components {
		c.Update(msg)
	}
}
func (l *Layout) GetView() tview.Primitive {
	return l.view
}

func NewLayout(rateLimit int, fname string) *Layout {
	container := tview.NewFlex().
		SetFullScreen(true).
		SetDirection(tview.FlexRow)

	header := tview.NewTable().
		SetFixed(2, 2).
		SetCell(0, 0, tview.NewTableCell("Log Source         ")).
		SetCell(0, 1, tview.NewTableCell(fname))
	container.AddItem(header, 2, 1, false)

	summary := tview.NewFlex().
		SetFullScreen(true).
		SetDirection(tview.FlexColumn)
	container.AddItem(summary, 0, 10, false)

	mostHits := &MostHits{
		view: tview.NewTable().SetFixed(6, 2),
	}
	summary.AddItem(mostHits.view, 0, 1, false)

	requestCounts := &RequestCounts{
		view: tview.NewTable().SetFixed(3, 2),
	}
	summary.AddItem(requestCounts.view, 0, 1, false)

	UniqueVisitors := &UniqueVisitors{
		view: tview.NewTable().SetFixed(2, 2),
	}
	summary.AddItem(UniqueVisitors.view, 0, 1, false)

	footer := &Footer{
		view:      tview.NewTable().SetFixed(1, 2),
		bucket:    make([]int, 12),
		rateLimit: rateLimit,
	}
	container.AddItem(footer.view, 1, 1, false)

	l := &Layout{
		view: container,
		components: []Component{
			mostHits,
			requestCounts,
			UniqueVisitors,
			footer,
		},
	}

	for _, c := range l.components {
		c.Init()
	}

	return l
}
