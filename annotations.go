package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func NewAnnotationsFromFile(path string) (map[string]string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	return NewAnnotations(f)
}

func annotationLineError(lineNumber int) error {
	return fmt.Errorf("line %d is malformed, was expecting line with format %q", lineNumber, "<key>=<quoted value>")
}

func NewAnnotations(f io.Reader) (map[string]string, error) {
	sc := bufio.NewScanner(f)
	annotations := map[string]string{}
	lineNumber := 0

	for sc.Scan() {
		line := sc.Text()
		lineNumber++

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		delimiterIndex := strings.Index(sc.Text(), "=")
		if delimiterIndex == -1 {
			return nil, annotationLineError(lineNumber)
		}

		k := line[0:delimiterIndex]
		v := line[delimiterIndex+1:]

		v, err := strconv.Unquote(v)
		if err != nil {
			return nil, annotationLineError(lineNumber)
		}

		annotations[k] = v
	}
	return annotations, nil
}
