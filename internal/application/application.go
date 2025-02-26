package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

type Request struct {
	Expression string `json:"expression"`
}

type Result struct {
	Result string `json:"result"`
}

type Error struct {
	Error string `json:"error"`
}

type Id struct {
	Id int `json:"id"`
}

var (
	taskIDCounter int
	mu            sync.Mutex
	tasks         []Task
	expIDCounter  int
	exp_mu        sync.Mutex
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func valid_check(expression string) bool {
	if regexp.MustCompile(`[a-zA-Zа-яА-Я]`).MatchString(expression) {
		return false
	}
	counter3 := 0
	counter4 := 0
	// проверяем, чтобы количество знаков операций скобок не превышало количества чисел. Иначе возвращаем фолс
	for i := 0; i < len(expression); i++ {
		if expression[i] == '0' || expression[i] == '1' || expression[i] == '2' || expression[i] == '3' || expression[i] == '4' || expression[i] == '5' || expression[i] == '6' || expression[i] == '7' || expression[i] == '8' || expression[i] == '9' {
			counter3 += 1
		} else if expression[i] == '*' || expression[i] == '/' || expression[i] == '-' || expression[i] == '+' {
			counter4 += 1
		}
	}
	if counter3 <= counter4 {
		return false
	}
	// проверяем, есть ли скобки в выражении. если да, то сверяем кол-во закрывающихся и открывающихся
	if strings.Contains(expression, "(") {
		counter1 := 0
		counter2 := 0
		for i := 0; i < len(expression); i++ {
			if expression[i] == '(' {
				counter1 += 1
			} else if expression[i] == ')' {
				counter2 += 1
			}
		}

		if counter1 != counter2 {
			return false
		}

	}
	return true
}

// разделение выражения без скобок чиста
func parseMultAndDiv(expression string) {
	err := godotenv.Load("iternal/application/.env")
	if err != nil {
		log.Printf(err.Error())
	}

	expression = strings.Replace(expression, "+", " + ", -1)
	expression = strings.Replace(expression, "-", " - ", -1)
	expression = strings.Replace(expression, "*", " * ", -1)
	expression = strings.Replace(expression, "/", " / ", -1)
	splited := strings.Fields(expression)

	for i := range len(splited) {

		if i+2 < len(splited) {
			if string(splited[i]) == "*" || string(splited[i]) == "/" && string(splited[i+2]) == "/" || string(splited[i+2]) == "*" {
				mu.Lock()
				taskIDCounter++
				id := taskIDCounter
				mu.Unlock()

				arg1, _ := strconv.Atoi(splited[i-1])
				arg2, _ := strconv.Atoi(splited[i+1])

				time := 1000
				if string(splited[i]) == "*" {
					time, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
				} else {
					time, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
				}

				task := Task{
					ID:            id,
					Arg1:          arg1,
					Arg2:          arg2,
					Operation:     splited[i],
					OperationTime: time,
				}
				fmt.Println(task)

				continue
				// task сразу отправляется к демону, после чего мы результат заменяем на исходное выражение
			}
		}
		if string(splited[i]) == "*" || string(splited[i]) == "/" {
			mu.Lock()
			taskIDCounter++
			id := taskIDCounter
			mu.Unlock()

			arg1, _ := strconv.Atoi(splited[i-1])
			arg2, _ := strconv.Atoi(splited[i+1])

			time := 1000
			if string(splited[i]) == "*" {
				time, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
			} else {
				time, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
			}

			task := Task{
				ID:            id,
				Arg1:          arg1,
				Arg2:          arg2,
				Operation:     splited[i],
				OperationTime: time,
			}
			tasks = append(tasks, task)
		}
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Критическая ошибка: %s", err.Error())
		if valid_check(request.Expression) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		exp_mu.Lock()
		json.NewEncoder(w).Encode(Id{Id: expIDCounter})
		exp_mu.Unlock()
	} else {
		log.Printf("Выражение успешно принято для вычислений")
		w.WriteHeader(http.StatusCreated)
		exp_mu.Lock()
		json.NewEncoder(w).Encode(Id{Id: expIDCounter})
		expIDCounter++
		exp_mu.Unlock()
	}
	count(request.Expression)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	task := tasks[0]
	tasks = tasks[1:]
	response := map[string]interface{}{
		"task": task,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (a *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", createHandler)
	http.HandleFunc("/internal/task", getTaskHandler)
	log.Printf("Сервер запущен")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
		return nil
	} else {
		return err
	}
}
