package main

import(
	"net"
	"io"
	"log"
	"bufio"
)

func advancedEcho(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if (err != nil) {
		log.Fatalln("Unable To Read")
	}
	log.Printf("Read %d bytes: %s", len(s), s)
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(s); err != nil {
		log.Fatalln("Unable To Write")
	}
	writer.Flush()
}

func echo(conn net.Conn) {
	defer conn.Close()
	b := make([] byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client Disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected Error")
			break
		}
		log.Printf("Received %d bytes: %s \n", size, string(b))
		resp := [] byte (" Fuck You ! ")
		if _, err := conn.Write(resp); err != nil {
			log.Fatalln("Unable to Write")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":2080")
	if err != nil {
		log.Fatalln("Unable to bind port")
	} 
	for {
		conn, err := listener.Accept()
		log.Println("Connection Received")
		if err != nil {
			log.Fatalln("Unable to accept connection", err)
		}
		go advancedEcho(conn)
	}
}