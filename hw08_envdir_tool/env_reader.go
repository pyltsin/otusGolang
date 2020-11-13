package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]string

const forbiddenSymbol = "="

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envDir, err := ioutil.ReadDir(dir)
	result := Environment{}
	if err != nil {
		return nil, fmt.Errorf("error while try to open dir: %w", err)
	}

	var path string
	for _, val := range envDir {
		path = filepath.Join(dir, val.Name())
		if !val.IsDir() {
			if strings.Contains(val.Name(), forbiddenSymbol) {
				continue
			}
			env, err := processFile(path)
			if err != nil {
				return nil, err
			}
			result[val.Name()] = env
		}
	}
	return result, nil
}

func processFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	fileStat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting stat: %w", err)
	}
	if fileStat.Size() == 0 {
		return "", nil
	}
	r := bufio.NewReader(file)
	envParam, _, err := r.ReadLine()

	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	result := strings.TrimRight(string(bytes.ReplaceAll(envParam, []byte("\x00"), []byte("\n"))), "\t \n")
	return result, nil
}
