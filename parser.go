package httptop

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Record struct {
	RemoteAddr string
	Userid     string
	Time       time.Time
	Method     string
	Request    string
	StatusCode int
	Bytes      int
}

func parseFields(s string) []string {
	pattern := `^(?P<client>\S+) \S+ (?P<userid>\S+) \[(?P<datetime>[^\]]+)\] "(?P<method>[A-Z]+) (?P<request>[^ "]+)? HTTP/[0-9.]+" (?P<status>[0-9]{3}) (?P<size>[0-9]+|-)`
	fields := regexp.MustCompile(pattern).FindStringSubmatch(s)

	return fields
}

func NewRecord(line string) (Record, error) {
	record := Record{}
	fields := parseFields(line)

	if len(fields) < 8 {
		msg := fmt.Sprintf("Missing some fields in %v", fields)
		return record, errors.New(msg)
	}

	record.RemoteAddr = fields[1]
	record.Userid = fields[2]

	// https://github.com/golang/go/issues/20797
	t, err := time.ParseInLocation("02/Jan/2006:15:04:05 -0700", fields[3], time.UTC)
	if err != nil {
		return record, err
	}
	record.Time = t

	record.Method = fields[4]
	record.Request = fields[5]

	statusCode, err := strconv.Atoi(fields[6])
	if err != nil {
		return record, err
	}
	record.StatusCode = statusCode

	bytes, err := strconv.Atoi(fields[7])
	if err != nil {
		return record, err
	}
	record.Bytes = bytes

	return record, nil
}
