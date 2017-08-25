package main

import (
	"os"
	"strings"
)

func main() {
	currentDirectory, _ := os.Getwd()

	if len(os.Args) < 2 {
		speak("error> missing inventory", true)
	}

	if len(os.Args) == 2 {
		os.Args = append(os.Args, "ls")
	}

	args := ""
	if len(os.Args) >= 3 {
		args = strings.Join(os.Args[3:], " ")
	}

	newMarvin(marvinFile(".", "marvin.yml"), currentDirectory, os.Args[1], os.Args[2], args)
}
