package httptop

import (
	"sort"
)

func getSection(s string) string {
	idx := len(s)
	for i, r := range s {
		if i > 0 && r == '/' {
			idx = i
			break
		}
	}

	return s[:idx]
}

type SectionHits struct {
	Section string
	Hits    int
}

type SecionHitsSlice []SectionHits

func (s SecionHitsSlice) Len() int           { return len(s) }
func (s SecionHitsSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SecionHitsSlice) Less(i, j int) bool { return s[i].Hits < s[j].Hits }

func GetTopHits(records []Record, limit int) []SectionHits {
	hitCounter := make(map[string]int)

	for _, record := range records {
		section := getSection(record.Request)

		hitCounter[section]++
	}

	topHits := make(SecionHitsSlice, len(hitCounter))

	i := 0
	for k, v := range hitCounter {
		topHits[i] = SectionHits{Section: k, Hits: v}
		i++
	}

	sort.Sort(sort.Reverse(topHits))

	idx := limit
	if len(topHits) < limit {
		idx = len(topHits)
	}

	return topHits[:idx]
}
