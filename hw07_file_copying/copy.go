package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = toFile.Close()
		_ = fromFile.Close()
	}()

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
	bar := pb.Full.Start64(limit)

	barReader := bar.NewProxyReader(fromFile)

	_, err = io.Copy(toFile, io.LimitReader(barReader, limit))

	bar.Finish()

	if err != nil {
		return err
	}
	return nil
}
