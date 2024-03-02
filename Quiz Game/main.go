package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question, answer string
}

func main() {
	// Get command line flags
	filename := flag.String("f", "problems.csv", "csv file to read problems from")
	timeLimit := flag.Int("t", 30, "time limit in seconds")
	flag.Parse()

	// Open CSV file
	file, err := os.Open(*filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Read CSV file
	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	problems := []Problem{}
	for _, line := range lines {
		problems = append(problems, Problem{line[0], line[1]})
	}

	numQuestions := len(problems)
	score := 0

	// Start timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Start looping through problems
problemLoop:
	for problemNum, problem := range problems {
		// Ask question and get response from user
		fmt.Printf("#%d: %s\n", problemNum+1, problem.question)

		answerChannel := make(chan string)
		go func() {
			var userAnswer string
			fmt.Scanln(&userAnswer)
			answerChannel <- userAnswer
		}()

		select {
		case <-timer.C:
			// Output user's score
			fmt.Println()
			fmt.Println("Time's up!")
			break problemLoop
		case userAnswer := <-answerChannel:
			// Update score
			if strings.Compare(userAnswer, problem.answer) == 0 {
				score++
			}
		}
	}
	// Output user's score
	percentage := float32(score) / float32(numQuestions) * 100
	fmt.Printf("Final score: %d out of %d (%.2f%%)\n", score, numQuestions, percentage)
}
