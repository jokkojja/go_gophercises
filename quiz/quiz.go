package main;
import "fmt"
import (
  "encoding/csv"
  "os"
  "strconv"
  "time"
)

type Problem struct {
  question string 
  answer int
}

type Problems struct {
  problems []Problem
}

func FromFile(file *os.File) Problems{
  var problems Problems = Problems{}

  reader := csv.NewReader(file)

  reader.Comma = ','
  reader.FieldsPerRecord = -1

  records, err := reader.ReadAll();
  
  if err != nil {
    panic(fmt.Errorf("Can't parse file. Err: %w", err))
  }
  for _, val := range records{
    ans, _ := strconv.ParseInt(val[1], 10, 8)
    problems.problems = append(problems.problems, Problem{
        question: val[0],
        answer: int(ans),
      })
  }

  return problems
}

func timer(secs_for_quiz int, ch chan bool) {
  timer := time.NewTimer(time.Duration(secs_for_quiz) * time.Second)
  <-timer.C
  ch <- true
}

func game(){
  
  ch := make(chan bool)
  answCh := make(chan int)

  file, err := os.Open("quiz/problems.csv")
  defer file.Close()
  if err != nil{
    panic(fmt.Errorf("Can't open file. Err: %w", err))
  }

  var problems Problems = FromFile(file);
  var total_right_answers int = 0;
  var total_wrong_answers int = 0;

  fmt.Println("Game started. Please, answer the questions");

  go timer(5, ch)
  
  for _, problem := range problems.problems{
    go func() {
      var user_answer int;
      fmt.Println(problem.question);
      fmt.Scanln(&user_answer)
      answCh <- user_answer;
    }()

    select {
    case is_game_finished := <-ch:
      if is_game_finished == true{
        fmt.Println("Time is up! U lost!")
        return
      }

    case user_answer := <- answCh:
      if user_answer == problem.answer {
        total_right_answers++;
      } else {
        total_wrong_answers++;
      }
    }
  }

    fmt.Println("Game finished. Total right answers: ", total_right_answers, " total wrong ansers: ", total_wrong_answers);
}

func main(){
  
  game()

}
