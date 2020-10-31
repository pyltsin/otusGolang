package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type BarReader struct {
	file *os.File
	pb   *pb.ProgressBar
}

func (receiver *BarReader) Read(p []byte) (n int, err error) {
	read, err := receiver.file.Read(p)
	if err == io.EOF {
		receiver.pb.Finish()
	}
	receiver.pb.Add(read)
	return read, err
}

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	stat, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get stat for source file: %w", err)
	}
	fileSize := stat.Size()
	if offset >= fileSize {
		return ErrOffsetExceedsFileSize
	}
	if fileSize == 0 {
		return ErrUnsupportedFile
	}
	if limit > fileSize-offset || limit == 0 {
		limit = fileSize - offset
	}

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return err
	}
	bar := pb.StartNew(int(limit))

	pbReader := BarReader{fromFile, bar}

	bufferSize := 2 << 10
	buffer := make([]byte, bufferSize)

	_, err = io.CopyBuffer(toFile, io.LimitReader(&pbReader, limit), buffer)

	if err != nil {
		return err
	}
	return nil
}
