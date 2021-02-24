package main

import (
	runner "github.com/alwashali/cmdlog/internal/runner"
)

func main() {

	options := runner.ParseOptions()
	r := runner.New(options)
	r.Execute()
}
