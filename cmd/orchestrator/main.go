package main

import "fourth/internal/application"

func main() {
	app := application.NewOrchestrator()
	app.Run()
}
