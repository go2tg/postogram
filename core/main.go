package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

const (
	HOST              = "localhost"
	PORT              = "2525"
	GREETINGS         = "220 localhost.local\n"
	EHLO              = "EHLO"
	R_EHLO            = "500 \n"
	HELO              = "HELO"
	R_HELO_DOMAIN     = "250 \n"
	MAIL_FROM         = "MAIL FROM:"
	R_MAIL_FROM       = "250 \n"
	RCPT              = "RCPT TO:"
	R_RCPT            = "250 \n"
	DATA              = "DATA"
	R_DATA            = "354 \n"
	R_CRLF_POINT_CRLF = "250 \n"
	QUIT              = "QUIT"
	R_QUIT            = "221 \n"
)

func main() {

	var w sync.WaitGroup

	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Port error", err.Error())
		return
	}

	defer l.Close()

	for i := 0; i < 3; i++ { // Запустить 3 коннекта
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Port error", err.Error())
			return
		}
		w.Add(1)
		go handleConnection(conn, &w) // при QUIT завершитья рутина
	}
	w.Wait() // Ждем завершения 3 рутин
}

func handleConnection(conn net.Conn, w *sync.WaitGroup) {

	conn.Write([]byte(GREETINGS))

	for {

		// ... CRLF = CR = 0x0d + LF = 0x0a
		netData, err := bufio.NewReader(conn).ReadBytes(0x0a)
		if err != nil {
			fmt.Println(err)
			return
		}

		t := strings.TrimSpace(string(netData))

		if strings.Contains(t, EHLO) || strings.Contains(t, "ehlo") {
			conn.Write([]byte(R_EHLO))
		}

		if strings.Contains(t, HELO) || strings.Contains(t, "helo") {
			conn.Write([]byte(R_HELO_DOMAIN))
		}

		if strings.Contains(t, MAIL_FROM) || strings.Contains(t, "mail FROM:") {
			fmt.Println("FROM MAIL: ", parsCommand(t))
			conn.Write([]byte(R_MAIL_FROM))
		}

		if strings.Contains(t, RCPT) || strings.Contains(t, "rcpt TO:") {
			fmt.Println("RCPT TO: ", parsCommand(t))
			conn.Write([]byte(R_RCPT))
		}

		if t == DATA || t == "data" {
			conn.Write([]byte(R_DATA))
		}

		if bytes.Contains(netData, []byte{0x2e, 0x0d, 0x0a}) {
			conn.Write([]byte(R_CRLF_POINT_CRLF))
			fmt.Printf("%x", t)
		}

		if t == QUIT || t == "quit" {
			conn.Write([]byte(R_QUIT))
			conn.Close()
			w.Done()
		}
	}
	conn.Close()
}

func parsCommand(t string) string {
	strings.Trim(t, ":")
	re := regexp.MustCompile("\\<(.*?)\\>")
	match := re.FindStringSubmatch(t)
	if match != nil {
		return match[1]
	} else {
		return ""
	}
}
