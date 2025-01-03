package main

import (
	"concernedmate/SurvivorGame/engines"
)

func main() {
	// start server
	go engines.OpenRoom("")
	engines.Server()
}
