package main

import 	"github.com/zakharkaverin1/calculator/internal/application"


func main() {
	app := application.NewOrchestrator()
	app.Run()
}
