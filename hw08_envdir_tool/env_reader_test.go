package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	baseDir := "testdata/env"
	countFiles := 4

	t.Run("read env vars", func(t *testing.T) {
		envs, err := ReadDir(baseDir)
		assert.Equal(t, nil, err)
		assert.Equal(t, countFiles, len(envs))
	})

	t.Run("skip file with =", func(t *testing.T) {
		invalidFileName := "test="
		tmpFile, _ := ioutil.TempFile(baseDir, invalidFileName)
		defer os.Remove(tmpFile.Name())
		envs, err := ReadDir(baseDir)
		assert.Equal(t, nil, err)
		assert.Equal(t, countFiles, len(envs))
	})

	t.Run("collect empty env string from empty file", func(t *testing.T) {
		fileName := "test"
		tmpFile, _ := ioutil.TempFile(baseDir, fileName)
		defer os.Remove(tmpFile.Name())
		envs, err := ReadDir(baseDir)
		assert.Equal(t, nil, err)
		assert.Equal(t, countFiles+1, len(envs))
		assert.Equal(t, "", envs[filepath.Base(tmpFile.Name())])
	})

	t.Run("remove whitespace characters from env var", func(t *testing.T) {
		fileName := "test"
		tmpFile, _ := ioutil.TempFile(baseDir, fileName)
		defer os.Remove(tmpFile.Name())
		_ = ioutil.WriteFile(tmpFile.Name(), []byte("first line\t \t \nsecond line"), os.ModePerm)
		envs, err := ReadDir(baseDir)
		assert.Equal(t, nil, err)
		assert.Equal(t, countFiles+1, len(envs))
		assert.Equal(t, "first line", envs[filepath.Base(tmpFile.Name())])
	})
}
