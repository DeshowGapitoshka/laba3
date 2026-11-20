package sort

import (
	"fmt"
	"lr3/parser"
	"sort"
	"strings"
	"time"
)

type Candidate struct {
	Dist       uint32
	LineNumber int
	Label      string
}

func Distance(line1 []uint32, line2 []uint32) uint32 {
	var i, j int
	var intersection uint32
	for i < len(line1) && j < len(line2) {

		if line1[i] == line2[j] {
			intersection++
			i++
			j++
			continue
		} else if line1[i] < line2[j] {
			i++
			continue
		} else {
			j++
			continue
		}
	}
	return uint32(len(line1)) + uint32(len(line2)) - 2*intersection
}

func Selection(parsed *parser.ParsedData, queryIndex int, query []uint32, train [][]uint32) (string, time.Duration) {
	var candidates []Candidate
	start := time.Now()
	for i := 0; i < len(train); i++ {
		dist := Distance(query, train[i])
		candidates = append(candidates, Candidate{
			Dist:       dist,
			LineNumber: parsed.Train[i].LineNumber,
			Label:      parsed.Train[i].Label,
		})
	}
	end := time.Since(start)
	if len(candidates) == 0 {
		return fmt.Sprintf("Q %d -", queryIndex+1), end
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Dist != candidates[j].Dist {
			return candidates[i].Dist < candidates[j].Dist
		}
		return candidates[i].LineNumber < candidates[j].LineNumber
	})
	k := parsed.K
	if k > len(candidates) {
		k = len(candidates)
	}
	neighbors := candidates[:k]

	votes := make(map[string]int)
	for _, n := range neighbors {
		votes[n.Label]++
	}

	bestLabel := ""
	bestCount := -1
	for _, n := range neighbors {
		if votes[n.Label] > bestCount {
			bestCount = votes[n.Label]
			bestLabel = n.Label
		}
	}

	var out strings.Builder
	fmt.Fprintf(&out, "Q %d %s", queryIndex+1, bestLabel)
	for _, n := range neighbors {
		fmt.Fprintf(&out, " %d:%d", n.LineNumber, n.Dist)
	}
	return out.String(), end
}
