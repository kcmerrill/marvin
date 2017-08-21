package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		speak("error> missing inventory", true)
	}

	if len(os.Args) == 2 {
		//speak("error> missing task, or at the very least, shell arguments", true)
		os.Args = append(os.Args, "ls")
	}

	args := ""
	if len(os.Args) >= 3 {
		args = strings.Join(os.Args[3:], " ")
	}

	newMarvin(marvinFile(".", "marvin.yml"), os.Args[1], os.Args[2], args)
}
