package persistence

import (
	"bufio"
	"os"
	"sync"
)

type AofLogger struct {
	File *os.File
	Mu   sync.Mutex
	Path string
}

func NewAofLogger(path string) (*AofLogger, error) {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &AofLogger{File: file}, nil

}

func (a *AofLogger) Write(command string) error {

	a.Mu.Lock()
	defer a.Mu.Unlock()

	_, err := a.File.WriteString(command + "\n")

	return err

}

func (a *AofLogger) Close() error {

	a.Mu.Lock()
	defer a.Mu.Unlock()

	err := a.File.Close()

	return err

}

func (a *AofLogger) Read() ([]string, error) {

	file, err := os.Open(a.Path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var commands []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		commands = append(commands, text)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}
