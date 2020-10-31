package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func compare(file1, file2 string) bool {
	b1, err := ioutil.ReadFile(file1)
	if err != nil {
		return false
	}

	b2, err := ioutil.ReadFile(file2)
	if err != nil {
		return false
	}

	return bytes.Equal(b1, b2)
}

type test struct {
	offset int64
	limit  int64
	answer string
}

func TestCopy(t *testing.T) {
	defer os.Remove("testdata/temp.txt")

	for _, tst := range [...]test{
		{
			offset: 0,
			limit:  0,
			answer: "testdata/out_offset0_limit0.txt",
		}, {
			offset: 0,
			limit:  10,
			answer: "testdata/out_offset0_limit10.txt",
		}, {
			offset: 0,
			limit:  1000,
			answer: "testdata/out_offset0_limit1000.txt",
		}, {
			offset: 0,
			limit:  10000,
			answer: "testdata/out_offset0_limit10000.txt",
		}, {
			offset: 100,
			limit:  1000,
			answer: "testdata/out_offset100_limit1000.txt",
		}, {
			offset: 6000,
			limit:  1000,
			answer: "testdata/out_offset6000_limit1000.txt",
		},
	} {
		err := Copy("testdata/input.txt", "testdata/temp.txt", tst.offset, tst.limit)
		require.Nil(t, err)
		isEqual := compare("testdata/temp.txt", tst.answer)
		require.True(t, isEqual)
	}
}
