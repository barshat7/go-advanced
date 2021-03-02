package readerwriter

import (
	"fmt"
	"os"
)

type FooReader struct{}

func (f *FooReader) Read(b []byte) (int, error) {
	fmt.Println("in >")
	return os.Stdin.Read(b)
}

type FooWriter struct {}

func (f *FooWriter) Write(b []byte) (int, error) {
	fmt.Println("out >")
	return os.Stdout.Write(b)
}

