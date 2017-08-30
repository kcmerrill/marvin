package main

import (
	"flag"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	args := flag.Args()

	currentDirectory, _ := os.Getwd()

	if len(args) < 2 {
		speak("error> missing inventory", true)
	}

	if len(args) == 2 {
		os.Args = append(os.Args, "ls")
	}

	passThrough := ""
	if len(args) >= 3 {
		passThrough = strings.Join(args[2:], " ")
	}
	newMarvin(marvinFile(".", "marvin.yml"), currentDirectory, os.Args[1], os.Args[2], passThrough)
}
