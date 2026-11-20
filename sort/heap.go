package sort

import (
	"fmt"
	"lr3/parser"
	"sort"
	"strings"
	"time"
)

type MaxHeap struct {
	data []Candidate
	k    int
}

func isWorse(a, b Candidate) bool {
	if a.Dist != b.Dist {
		return a.Dist > b.Dist
	}
	return a.LineNumber > b.LineNumber
}

func (h *MaxHeap) siftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if isWorse(h.data[index], h.data[parent]) == false {
			break
		}
		h.data[index], h.data[parent] = h.data[parent], h.data[index]
		index = parent
	}
}

func (h *MaxHeap) siftDown(index int) {
	n := len(h.data)
	for {
		left := 2*index + 1
		right := 2*index + 2
		largest := index
		if left < n && isWorse(h.data[left], h.data[largest]) {
			largest = left
		}
		if right < n && isWorse(h.data[right], h.data[largest]) {
			largest = right
		}
		if largest == index {
			break
		}
		h.data[index], h.data[largest] = h.data[largest], h.data[index]
		index = largest
	}
}

func (h *MaxHeap) Push(c Candidate) {
	if len(h.data) < h.k {
		h.data = append(h.data, c)
		h.siftUp(len(h.data) - 1)
	} else {
		if isWorse(h.data[0], c) {
			h.data[0] = c
			h.siftDown(0)
		}
	}
}

func (h *MaxHeap) GetSorted() []Candidate {
	result := make([]Candidate, len(h.data))
	copy(result, h.data)
	// Сортируем по возрастанию (лучшие — первыми)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Dist != result[j].Dist {
			return result[i].Dist < result[j].Dist
		}
		return result[i].LineNumber < result[j].LineNumber
	})
	return result
}

func SelectionHeap(parsed *parser.ParsedData, queryIndex int, query []uint32, train [][]uint32) (string, time.Duration) {
	k := parsed.K
	if k > len(train) {
		k = len(train)
	}
	heap := &MaxHeap{k: k}
	start := time.Now()
	for i := 0; i < len(train); i++ {
		dist := Distance(query, train[i])
		c := Candidate{
			Dist:       dist,
			LineNumber: parsed.Train[i].LineNumber,
			Label:      parsed.Train[i].Label,
		}
		heap.Push(c)
	}
	end := time.Since(start)
	if len(heap.data) == 0 {
		return fmt.Sprintf("Q %d -", queryIndex+1), end
	}

	neighbors := heap.GetSorted()

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
