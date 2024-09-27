package utils

import (
	"bufio"
	"io"
	"os"
)

func ReadLineFromFile(path string, lineNum int) (line string, err error) {
	readFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	line, _, err = ReadLine(readFile, lineNum)
	return line, err
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}
