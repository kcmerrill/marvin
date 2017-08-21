package main

import (
	"os"

	"github.com/fatih/color"
)

func speak(msg string, exit bool) {
	if exit {
		color.Red(msg)
		os.Exit(42)
	} else {
		color.Green(msg)
	}
}
