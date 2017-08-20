package main

import (
	"fmt"
	"os"
)

func speak(msg string, exit bool) {
	fmt.Println(msg)
	if exit {
		os.Exit(42)
	}
}
