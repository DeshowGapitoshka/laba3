package sort

import (
	"fmt"
	"lr3/parser"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func less(a, b Candidate) bool {
	if a.Dist != b.Dist {
		return a.Dist < b.Dist
	}
	return a.LineNumber < b.LineNumber
}

func partition(arr []Candidate, left, right int) int {
	pivotIndex := left + rand.Intn(right-left+1)
	arr[pivotIndex], arr[right] = arr[right], arr[pivotIndex]
	pivot := arr[right]

	i := left
	for j := left; j < right; j++ {
		if less(arr[j], pivot) {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[right] = arr[right], arr[i]
	return i
}

func quickSelect(arr []Candidate, k int) {
	if k <= 0 || k >= len(arr) {
		return
	}
	left, right := 0, len(arr)-1
	target := k - 1
	for left <= right {
		pivotIndex := partition(arr, left, right)
		if pivotIndex == target {
			break
		} else if pivotIndex < target {
			left = pivotIndex + 1
		} else {
			right = pivotIndex - 1
		}
	}
}

func SelectionQuick(parsed *parser.ParsedData, queryIndex int, query []uint32, train [][]uint32) (string, time.Duration) {
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
	k := parsed.K
	if k > len(candidates) {
		k = len(candidates)
	}

	quickSelect(candidates, k)

	neighbors := candidates[:k]
	sort.Slice(neighbors, func(i, j int) bool {
		return less(neighbors[i], neighbors[j])
	})

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
