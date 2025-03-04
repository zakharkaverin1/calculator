package application

import (
	"bytes"
	"encoding/json"
	"fourth/pkg/calculation"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Id  int     `json:"id"`
	Res float64 `json:"res"`
}

type Agent struct {
	power int
	url   string
}

func NewAgent() *Agent {
	p, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		p = 1
	}
	return &Agent{power: p, url: "http://:8080"}
}

func (a *Agent) Run() {
	for i := 0; i < a.power; i++ {
		log.Printf("Начинаем работу демона номер %d", i)
		go a.worker(i)
	}
	select {}
}

func (a *Agent) worker(id int) {
	for {
		resp, err := http.Get(a.url + "/internal/task")
		if err != nil {
			log.Printf("1 демон%d ошибка: %v", id, err)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(1 * time.Second)
			continue
		}
		var taksa struct {
			Task struct {
				ID            string  `json:"id"`
				Arg1          float64 `json:"arg1"`
				Arg2          float64 `json:"arg2"`
				Operation     string  `json:"operation"`
				OperationTime int     `json:"operation_time"`
			} `json:"task"`
		}
		err = json.NewDecoder(resp.Body).Decode(&taksa)
		if err != nil {
			log.Printf("2 демон%d ошибка: %v", id, err)
			continue
		}
		time.Sleep(time.Duration(taksa.Task.OperationTime) * time.Millisecond)
		result, err := calculation.Compute(taksa.Task.Operation, taksa.Task.Arg1, taksa.Task.Arg2)
		if err != nil {
			log.Printf("3 демон%d ошибка: %v", id, err)
			continue
		}
		idshka, _ := strconv.Atoi(taksa.Task.ID)
		response := Response{Id: idshka, Res: result}
		jsonResp, _ := json.Marshal(response)
		respPost, _ := http.Post(a.url+"/internal/task", "application/json", bytes.NewReader(jsonResp))

		resp.Body.Close()
		if respPost.StatusCode != http.StatusOK {
			log.Printf("Ошибка запроса")
		} else {
			log.Printf("Демон выполнил свою работу")
		}
	}
}
