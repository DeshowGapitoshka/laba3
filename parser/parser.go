package parser

import (
	"bytes"
	"fmt"
	"lr3/tokenizer"
	"strconv"
)

type TrainItem struct {
	Label      string
	Words      [][]byte
	LineNumber int
}

type QueryItem struct {
	Words      [][]byte
	LineNumber int
}

type ParsedData struct {
	K       int
	Train   []TrainItem
	Queries []QueryItem
}

func Parser(lines [][]byte) (*ParsedData, error, []string) {
	var errors []string
	data := &ParsedData{}
	const (
		header = iota
		train
		query
	)
	state := header
	flag := 0
	for i, line := range lines {
		lineNumber := i + 1
		if len(line) == 0 {
			continue
		}
		switch state {
		case header:
			if flag == 0 {
				fields := bytes.Fields(line)
				if len(fields) != 2 || string(fields[0]) != "K" {
					return nil, fmt.Errorf("ERROR %d bad_header", lineNumber), errors
				}
				k, err := strconv.Atoi(string(fields[1]))
				if err != nil || k < 1 || k > 1_000_000_000 {
					return nil, fmt.Errorf("ERROR %d bad_header", lineNumber), errors
				}
				data.K = k
				flag = 1
			} else if flag == 1 {
				if string(line) != "TRAIN" {
					return nil, fmt.Errorf("ERROR %d bad_header", lineNumber), errors
				}
				state = train
			}
		case train:
			if string(line) == "QUERY" {
				state = query
				continue
			}

			ind := bytes.Index(line, []byte("||"))
			if ind == -1 {
				errors = append(errors, fmt.Sprintf("ERROR %d bad_train_line\n", lineNumber))
				continue
			}
			label := bytes.TrimSpace(line[:ind])
			text := line[ind+2:]
			if len(label) == 0 || bytes.ContainsAny(label, " \t") {
				errors = append(errors, fmt.Sprintf("ERROR %d bad_train_line\n", lineNumber))
				continue
			}
			words := tokenizer.Tokenizing(text)
			data.Train = append(data.Train, TrainItem{
				Label:      string(label),
				Words:      words,
				LineNumber: lineNumber,
			})
		case query:
			words := tokenizer.Tokenizing(line)
			data.Queries = append(data.Queries, QueryItem{
				Words:      words,
				LineNumber: lineNumber,
			})
		}
	}
	if flag != 1 || state == header {
		lastLineNumber := len(lines) + 1
		return nil, fmt.Errorf("ERROR %d bad_header", lastLineNumber), errors
	}
	return data, nil, errors
}
