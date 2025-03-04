package main

import (
	"fourth/internal/application"
	"log"
)

func main() {
	agent := application.NewAgent()
	log.Println("Агент запущен")
	agent.Run()
}
