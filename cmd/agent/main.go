package main

import (
	"log"

	"github.com/bulbosaur/calculator-with-authorization/config"
	agent "github.com/bulbosaur/calculator-with-authorization/internal/agent/service"
)

func main() {
	config.Init()

	log.Println("starting agent")
	agent.RunAgent()
}
