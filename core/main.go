package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	HOST              = "localhost"
	PORT              = "2525"
	GREETINGS         = "220\n"
	EHLO              = "EHLO"
	R_EHLO            = "500\n"
	HELO              = "HELO"
	R_DOMAIN          = "250\n"
	MAIL_FROM         = "MAIL FROM:"
	R_                = "250\n"
	RCPT              = "RCPT TO:"
	R_RCPT            = "250\n"
	DATA              = "DATA"
	R_DATA            = "354\n"
	QUIT              = "QUIT"
	R_CRLF_POINT_CRLF = "250\n"
	R_QUIT            = "221\n"
)

var (
	GREETINGS_COUNT = 1
)

func main() {

	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Port error", err.Error())
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Port error", err.Error())
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	for {

		if GREETINGS_COUNT == 1 {
			conn.Write([]byte(GREETINGS))
			GREETINGS_COUNT++
		}

		// ... CRLF = CR = 0x0d + LF = 0x0a
		netData, err := bufio.NewReader(conn).ReadBytes(0x0a)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s", netData)

		t := strings.TrimSpace(string(netData))

		if strings.Contains(t, EHLO) {
			conn.Write([]byte(R_EHLO))
		}

		if strings.Contains(t, HELO) {
			conn.Write([]byte(R_DOMAIN))
		}

		if strings.Contains(t, MAIL_FROM) {
			fmt.Println("FROM MAIL: ", parsCommand(t))
			conn.Write([]byte(R_))
		}

		if strings.Contains(t, RCPT) {
			fmt.Println("RCPT FROM: ", parsCommand(t))
			conn.Write([]byte(R_RCPT))
		}

		if t == DATA {
			conn.Write([]byte(R_DATA))
		}

		if bytes.Contains(netData, []byte{0x2e, 0x0d, 0x0a}) {
			conn.Write([]byte(R_CRLF_POINT_CRLF))
		}

		if t == QUIT {
			conn.Write([]byte(R_QUIT))
			conn.Close()
		}
	}

	conn.Close()
}

func parsCommand(t string) string {
	strings.Trim(t, ":")
	re := regexp.MustCompile("\\<(.*?)\\>")
	match := re.FindStringSubmatch(t)
	return match[1]
}
