package main

import (
	"log"

	runner "github.com/alwashali/cmdlog/internal/runner"
)

func main() {

	options := runner.ParseOptions()
	r, err := runner.New(options)
	if err != nil {
		log.Println("Error creating runner", err)
	}

	r.Execute()
}
