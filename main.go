package main

import (
	"bufio"
	"fmt"
	"log"
	"lr3/hash"
	"lr3/parser"
	sorting "lr3/sort"
	"lr3/tokenizer"
	"os"
	"sort"
	"strconv"
	"time"
)

type Indices struct {
	Train [][]uint32
	Query [][]uint32
}

func main() {
	args := os.Args
	if len(args) == 5 {
		driver(args, 1)
	} else if len(args) == 7 && args[5] == "--count" {
		count, err := strconv.Atoi(args[6])
		if err != nil {
			log.Fatal(err)
		}
		driver(args, 3)
		times, distanceTime := driver(args, count)
		if len(times) != 0 && len(distanceTime) != 0 {
			time := mediana(times)
			distanceT := mediana(distanceTime)
			fmt.Print("Время работы всего алгоритма: ")
			fmt.Println(time)
			fmt.Print("Время подсчёта расстояний: ")
			fmt.Println(distanceT)
		}
	} else {
		os.Exit(2)
	}

}

func driver(args []string, count int) ([]int, []int) {
	var times []int
	var distanceTime []int
	if args[1] == "--method" {
		var mseconds time.Duration
		for i := 0; i < count; i++ {
			var temp time.Duration
			var totalTemp time.Duration
			start := time.Now()
			outFile, err := os.Create(args[4])
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()
			w := bufio.NewWriter(outFile)

			data, err := os.ReadFile(args[3])
			var number uint32
			var indices Indices
			number = 1
			if err != nil {
				log.Println(err)
				return times, distanceTime
			}
			lines := tokenizer.SplitLines(data)
			parsed, err, errors := parser.Parser(lines)
			if err != nil {
				fmt.Fprintln(w, err)
				w.Flush()
				return times, distanceTime
			}
			table := hash.New()

			for i := 0; i < len(parsed.Train); i++ {
				var line []uint32
				for j := 0; j < len(parsed.Train[i].Words); j++ {
					wordStr := string(parsed.Train[i].Words[j])
					index, ok := table.Get(wordStr)

					if ok == false {
						index = number
						table.Insert(wordStr, number)
						number += 1
					}
					line = append(line, index)
				}
				sort.Slice(line, func(i, j int) bool {
					return line[i] < line[j]
				})
				if len(line) > 1 {
					unique := line[:1]
					for _, v := range line[1:] {
						if v != unique[len(unique)-1] {
							unique = append(unique, v)
						}
					}
					line = unique
				}
				indices.Train = append(indices.Train, line)
			}

			for i := 0; i < len(parsed.Queries); i++ {
				var line []uint32
				for j := 0; j < len(parsed.Queries[i].Words); j++ {
					wordStr := string(parsed.Queries[i].Words[j])
					index, ok := table.Get(wordStr)

					if ok == false {
						index = number
						table.Insert(wordStr, number)
						number += 1
					}
					line = append(line, index)
				}
				sort.Slice(line, func(i, j int) bool {
					return line[i] < line[j]
				})
				if len(line) > 1 {
					unique := line[:1]
					for _, v := range line[1:] {
						if v != unique[len(unique)-1] {
							unique = append(unique, v)
						}
					}
					line = unique
				}
				indices.Query = append(indices.Query, line)
			}
			var results []string
			for q := 0; q < len(parsed.Queries); q++ {
				var res string
				switch args[2] {
				case "sort":
					res, temp = sorting.Selection(parsed, q, indices.Query[q], indices.Train)
					totalTemp += temp
				case "heap":
					res, temp = sorting.SelectionHeap(parsed, q, indices.Query[q], indices.Train)
					totalTemp += temp
				case "select":
					res, temp = sorting.SelectionQuick(parsed, q, indices.Query[q], indices.Train)
					totalTemp += temp
				default:
					fmt.Printf("Неизвестный метод: %s\n", args[2])
					os.Exit(2)
				}
				results = append(results, res)
			}
			for _, err := range errors {
				fmt.Fprintln(w, err)
			}

			for _, r := range results {
				fmt.Fprintln(w, r)
			}
			mseconds = time.Since(start)
			w.Flush()
			times = append(times, int(mseconds.Nanoseconds()))
			distanceTime = append(distanceTime, int(totalTemp.Microseconds()))
		}

	} else {
		log.Println("Нет")
		os.Exit(2)
	}
	return times, distanceTime
}

func mediana(times []int) time.Duration {
	n := len(times)
	sort.Ints(times)
	if n == 1 {
		return time.Duration(times[0])
	}
	if n%2 == 0 {
		return time.Duration((times[n/2-1] + times[n/2]) / 2)
	} else {
		return time.Duration(times[n/2])
	}
}
