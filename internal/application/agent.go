package application

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Agent struct {
	power int
	url   string
}

type Response struct{
	id int `json:"id"`
	res float64 `json:"res"` 
}

func NewAgent() *Agent {
	p, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		p = 1
	}
	return &Agent{power: p, url: "http://localhost:8080"}
}

func (a *Agent) Run() {
	for i := 0; i < a.power; i++ {
		log.Printf("Начинаем работу демона номер", i)
		go a.worker(i)
	}
	select {}
}

func (a *Agent) worker(id int) {
	for {
		resp, err := http.Get(a.url + "/internal/task")
		defer resp.Body.Close()
		if err != nil {
			log.Printf("демон%d ошибка: %v", id, err)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(1 * time.Second)
			continue
		}
		var taksa Task
		err = json.NewDecoder(resp.Body).Decode(&taksa)
		if err != nil {
			log.Printf("демон%d ошибка: %v", id, err)
			continue
		} 
		time.Sleep(time.Duration(taksa.OperationTime) * time.Millisecond)
		
	}
}
