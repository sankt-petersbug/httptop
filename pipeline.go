package httptop

import (
	"time"

	"github.com/hpcloud/tail"
)

type Message struct {
	Records   []Record
	CreatedAt time.Time
}

func ReadFile(fname string) (*tail.Tail, error) {
	return tail.TailFile(fname, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		Logger:   tail.DiscardingLogger,
	})
}

func BatchRead(in <-chan Record, interval time.Duration) <-chan *Message {
	var buf []Record
	tick := time.Tick(interval)
	out := make(chan *Message)

	go func() {
		for {
			select {
			case now := <-tick:
				msg := &Message{Records: buf, CreatedAt: now}
				out <- msg
				buf = nil
			case record := <-in:
				buf = append(buf, record)
			}
		}
		close(out)
	}()
	return out
}

func ToRecord(in <-chan *tail.Line) <-chan Record {
	out := make(chan Record)
	go func() {
		for line := range in {
			record, err := NewRecord(line.Text)

			if err != nil {
				// TODO: log or alert user
			}

			out <- record
		}
		close(out)
	}()
	return out
}
