package main

import (
	"fmt"
	"franky/runner"
	"log"
	"time"
)

func main() {
	runner := runner.NewRunner("./test_app")

	fmt.Print("Starting pipeline execution")

	s := time.Now()

	if err := runner.Execute(); err != nil {
		log.Panic("Running pipeline failed")
	}

	fmt.Printf("Complete pipeline in %.2f seconds", time.Since(s).Seconds())
}
