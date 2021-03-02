package transactionlog

import (
	"os"
	"fmt"
	"bufio"
)

type Event struct {
	Sequence uint64
	EventType byte
	Key string
	Value string
}

const (
	EventDelete byte = 1
	EventPut byte = 2
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
    ReadEvents() (<-chan Event, <-chan error)
    Run()
}

type FileTransactionLogger struct {
	events chan <- Event
	errors <- chan error
	lastSequence uint64
	file *os.File
}

func (f *FileTransactionLogger) WriteDelete(key string) {
	f.events <- Event{EventType: EventDelete, Key: key}
}

func (f *FileTransactionLogger) WritePut(key, value string) {
	f.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (f *FileTransactionLogger) Err() <- chan error {
	return f.errors
}

func (f *FileTransactionLogger) Run() {
	fmt.Println("Running Transaction Logger")
	events := make(chan Event, 16)
	f.events = events
	errors := make(chan error, 1)
	f.errors = errors
	go func() {
		for e := range events {
			f.lastSequence++
			_, err := fmt.Fprintf(
				f.file,
				"%d\t%d\t%s\t%s\t\n",
				f.lastSequence, e.EventType, e.Key, e.Value,
			)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (f *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(f.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s\t\n",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			if f.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction sequence out of sequence")
				return
			}
		
			f.lastSequence = e.Sequence
			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
            outError <- fmt.Errorf("transaction log read failure: %w", err)
        }
	}()
	return outEvent, outError
}

/** Constructor */

func New(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
    if err != nil {
        return nil, fmt.Errorf("cannot open transaction log file: %w", err)
    }
    return &FileTransactionLogger{file: file}, nil
}