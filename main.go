package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file in the format of questions,answers")
	timeLimit := flag.Int("limit", 30, "time limit for quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file : %s", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Parsing error")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
problemLoop:
	for index, problem := range problems {
		fmt.Printf("Problem no %d -> %s : ", index+1, problem.question)
		ansChan := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s\n", &ans)
			ansChan <- ans
		}()
		select {
		case <-timer.C:
			{
				fmt.Println("\n----------Time Over----------")
				break problemLoop
			}
		case ans := <-ansChan:
			{
				if ans == problem.answer {
					correct++
				}
			}
		}
	}

	fmt.Printf("You scored %d out of %d", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
