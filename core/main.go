package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	HOST = "localhost"
	PORT = "2525"
)

func main() {

	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Port error", err.Error())
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Port error", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	//buf := make([]byte, 128)

	for {

		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.TrimSpace(netData) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Println("->", netData)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		_, err = conn.Write([]byte(myTime))
		if err != nil {
			fmt.Println("Error write TCP")
		}

	}

}

//func write(conn net.Conn, buf []byte) {
//	_, err := conn.Write(buf)
//	if err != nil {
//		fmt.Println("Error reading", err.Error())
//	}
//	//fmt.Print(buf)
//}

//func read(conn net.Conn, buf []byte) {
//	_, err := conn.Read(buf)
//	if err != nil {
//		fmt.Println("Error reading", err.Error())
//	} else {
//		fmt.Println(buf)
//	}
//}
