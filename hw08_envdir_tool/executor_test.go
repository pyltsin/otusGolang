package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("success run cmd", func(t *testing.T) {
		tmpCmd := []string{"echo"}
		code := RunCmd(tmpCmd, nil)
		assert.Equal(t, successCode, code)
	})
	t.Run("success set env variable", func(t *testing.T) {
		tmpCmd := []string{"echo"}
		tmpEnv := Environment{"TEST": "test"}
		code := RunCmd(tmpCmd, tmpEnv)
		assert.Equal(t, successCode, code)
		assert.Equal(t, "test", os.ExpandEnv("$TEST"))
	})
	t.Run("return code from command", func(t *testing.T) {
		tmpCmd := []string{"/bin/bash", "wrong_name"}
		code := RunCmd(tmpCmd, nil)
		assert.Equal(t, 127, code)
	})
	t.Run("unset env var", func(t *testing.T) {
		os.Setenv("TEST", "test")
		tmpCmd := []string{"echo"}
		tmpEnv := Environment{"TEST": ""}
		code := RunCmd(tmpCmd, tmpEnv)
		assert.Equal(t, successCode, code)
		remEnv, isPresent := os.LookupEnv("$TEST")
		assert.Equal(t, false, isPresent)
		assert.Equal(t, "", remEnv)
	})
}
