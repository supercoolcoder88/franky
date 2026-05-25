package main

import (
	"encoding/json"
	"fmt"
	"franky/runner"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic("failed to get working directory: ", err)
	}

	f, err := os.ReadFile(filepath.Join(wd, "franky.json"))
	if err != nil {
		log.Panic("failed to read config file: ", err)
	}

	var config runner.Config
	if err := json.Unmarshal(f, &config); err != nil {
		log.Panic("failed to unmarshal config file: ", err)
	}

	r := runner.NewRunner(&config)

	fmt.Println("Starting pipeline execution")

	s := time.Now()

	if err := r.Execute(); err != nil {
		log.Panic("Running pipeline failed: ", err)
	}

	fmt.Printf("Complete pipeline in %.2f seconds\n", time.Since(s).Seconds())
}
