package main

import (
	"log"

	"github.com/bulbosaur/web-calculator-golang/config"
	agent "github.com/bulbosaur/web-calculator-golang/internal/agent/service"
)

func main() {
	config.Init()

	log.Println("starting agent")
	agent.RunAgent()
}
