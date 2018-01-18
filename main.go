package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	filename := flag.String("file", "", "Path to transactions CSV. If empty, will use stdin")
	flag.Parse()

	r, err := getReader(*filename)
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}
	defer r.Close()
}

func getReader(file string) (*os.File, error) {
	if file == "" {
		return os.Stdin, nil
	}

	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return r, nil
}
