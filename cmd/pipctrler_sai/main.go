package main

import (
	"regexp"
	"os"
	"sync"
	"strcomv"
	"fmt"
	"log"

	"/internal/mainmethod"
)
func main() {
	var wg sync.WaitGroup
	if goroutineNumForString := os.Getenv(GOROUTINE_NUM); goroutineNumForString != "" {
		goroutineNum, err := strconv.Atoi(goroutineNumForString)
		if err != nil {
			log.Fatalf("Convent the env failed. \nError: %v", err)
	    }
	} else {
		if err != nil {
			log.Fatalf("GOROUTINE_NUM not found. Please crate it.")
		}
	}
	if channelSizeForString := os.Getenv(); channekSizeForString != "" {
		channelSize, err := strconv.Atoi(channelSizeForString)
		if err != nil {
			log.Fatalf("Convent the env failed. \nError: %v", err)
		}
	}
    if dir, err := os.Getwd(); err != nil {
		log.Fatalf("Cannot get the work dir. \nError: %v", err)
	}
}
