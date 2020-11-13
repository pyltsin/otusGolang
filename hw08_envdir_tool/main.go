package main

import (
	"log"
	"os"
)

func main() {
	inputArgs := os.Args

	if len(inputArgs) <= 3 {
		log.Fatal(`Expected: go-envdir /path/to/env/dir command arg1 arg2`)
	}

	pathToEnvDir := inputArgs[1]
	commands := inputArgs[2:]

	environmentParams, err := ReadDir(pathToEnvDir)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(RunCmd(commands, environmentParams))
}
