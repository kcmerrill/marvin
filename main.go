package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		speak("Invalid inventory or task", true)
	}

	args := ""
	if len(os.Args) >= 3 {
		args = strings.Join(os.Args[3:], " ")
	}

	newMarvin(marvinFile(".", "marvin.yml"), os.Args[1], os.Args[2], args)
}
