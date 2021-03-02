package main

import (
	_ "tcpscanner/simplescanner"
	rw "tcpscanner/readerwriter"
	"log"
	"fmt"
)
func main () {
	var reader rw.FooReader
	var writer rw.FooWriter

	input := make([]byte, 4096)

	s, err := reader.Read(input)

	if err != nil {
		log.Fatalln("Unable To Read")
	}

	fmt.Printf("Read %d bytes from stdin\n", s)

	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write")
	}
	fmt.Printf("Wrote %d bytes to stdout\n", s)
}