package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func marvinFile(dir, configFile string) string {
	currentDirectory, _ := os.Getwd()

	// make sure directory exists
	if _, dirExists := os.Stat(dir); dirExists != nil {
		dir = "."
	}

	// switch to the directory
	os.Chdir(dir)

	// keep on going on up
	for {
		if _, configExists := os.Stat(configFile); configExists == nil {
			marvinFile, marvinFileError := ioutil.ReadFile(configFile)
			if marvinFileError != nil {
				speak("Cannot read "+configFile, true)
			}
			return string(marvinFile)
		}
		curDir, curDirError := os.Getwd()
		if curDirError != nil {
			// huh, never seen this before
			speak("Unable to get directory information", true)
		}

		// finished yet?
		if abs, _ := filepath.Abs(curDir); abs == "/" {
			os.Chdir(currentDirectory)
			break
		}

		// go up a dir
		os.Chdir("..")
	}

	// nothing found
	return ""
}
