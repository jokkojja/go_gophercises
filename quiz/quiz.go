package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Problem struct {
	question string
	answer   int
}

type Problems struct {
	problems []Problem
}

func fromFile(file *os.File) Problems {
	var problems Problems

	reader := csv.NewReader(file)

	reader.Comma = ','
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("Can't parse file. Err: %w", err))
	}

	for _, val := range records {
		ans, _ := strconv.ParseInt(val[1], 10, 8)
		problems.problems = append(problems.problems, Problem{
			question: val[0],
			answer:   int(ans),
		})
	}

	return problems
}

func timer(secsForQuiz int, ch chan bool) {
	timer := time.NewTimer(time.Duration(secsForQuiz) * time.Second)
	<-timer.C
	ch <- true
}

func prepareRules() (filePath *string, secsForQuiz *int) {
	filePath = flag.String("filePath", "problems.csv", "Path to file with quiz")
	secsForQuiz = flag.Int("secsForQuiz", 30, "Secs for quizing")
	flag.Parse()

  return filePath, secsForQuiz
}

func game(filePath string, secsForQuiz int) {
	ch := make(chan bool)
	answCh := make(chan int)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		panic(fmt.Errorf("Can't open file. Err: %w", err))
	}

	var problems Problems = fromFile(file)
	var totalRightAnswers int
	var totalWrongAnswers int

	fmt.Println("Game started. Please, answer the questions")

	go timer(secsForQuiz, ch)

	for _, problem := range problems.problems {
		go func(problem Problem) {
			var userAnswer int
			fmt.Println(problem.question)
			fmt.Scanln(&userAnswer)
			answCh <- userAnswer
		}(problem)

		select {
		case isGameFinished := <-ch:
			if isGameFinished {
				fmt.Printf("Time (%v seconds) is up! You lost!", secsForQuiz)
				return
			}

		case userAnswer := <-answCh:
			if userAnswer == problem.answer {
				totalRightAnswers++
			} else {
				totalWrongAnswers++
			}
		}
	}

	fmt.Println("Game finished. Total right answers:", totalRightAnswers, "Total wrong answers:", totalWrongAnswers)
}

func main() {
  filePath, secsForQuiz := prepareRules()
	game(*filePath, *secsForQuiz)
}
