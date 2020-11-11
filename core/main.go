package main

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	HOST              = "localhost"
	PORT              = "2525"
	GREETINGS         = "220 localhost.local\n"
	EHLO              = "EHLO"
	EHLO_PY           = "ehlo"
	R_EHLO            = "500 \n"
	HELO              = "HELO"
	HELO_PY           = "helo"
	R_HELO_DOMAIN     = "250 \n"
	MAIL_FROM         = "MAIL FROM:"
	MAIL_FROM_PY      = "mail FROM:"
	R_MAIL_FROM       = "250 \n"
	RCPT              = "RCPT TO:"
	RCPT_PY           = "rcpt TO:"
	R_RCPT            = "250 \n"
	DATA              = "DATA"
	DATA_PY           = "data"
	R_DATA            = "354 \n"
	R_CRLF_POINT_CRLF = "250 \n"
	QUIT              = "QUIT"
	R_QUIT            = "221 \n"
)

var CRLF = []byte{0x0d, 0x0a, 0x2e, 0x0d, 0x0a}

func main() {

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
		go handleConnection(conn) // при QUIT завершитья рутина
	}
}

func handleConnection(conn net.Conn) {

	conn.Write([]byte(GREETINGS))
	buf := make([]byte, 128)

	for {

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error read TCP: ", err)
		}

		fmt.Printf("%s", buf)

		t := strings.TrimSpace(string(buf))

		// EHLO
		if bytes.Equal([]byte(EHLO), buf[0:4]) || bytes.Equal([]byte(EHLO_PY), buf[0:4]) {
			fmt.Println("EHLO - OK")
			conn.Write([]byte(R_EHLO))
		}

		// HELO
		if bytes.Equal([]byte(HELO), buf[0:4]) || bytes.Equal([]byte(HELO_PY), buf[0:4]) {
			fmt.Println("HELO - OK")
			conn.Write([]byte(R_HELO_DOMAIN))
		}

		// MAIL FROM:
		if bytes.Equal([]byte(MAIL_FROM), buf[0:10]) || bytes.Equal([]byte(MAIL_FROM_PY), buf[0:10]) {
			fmt.Println("MAIL FROM: - OK ", parsCommand(t))
			conn.Write([]byte(R_MAIL_FROM))
		}

		// RCPT TO:
		if bytes.Equal([]byte(RCPT), buf[0:8]) || bytes.Equal([]byte(RCPT_PY), buf[0:8]) {
			fmt.Println("RCPT TO: - OK ", parsCommand(t))
			conn.Write([]byte(R_RCPT))
		}

		// DATA
		if bytes.Equal([]byte(DATA), buf[0:4]) || bytes.Equal([]byte(DATA_PY), buf[0:4]) {
			fmt.Println("DATA - OK")
			conn.Write([]byte(R_DATA))
		}

		if bytes.Contains(buf, CRLF) {
			conn.Write([]byte(R_CRLF_POINT_CRLF))
			fmt.Printf("%x", t)
		}

		if strings.Contains(t, QUIT) || strings.Contains(t, "quit:") {
			conn.Write([]byte(R_QUIT))
			conn.Close()
			return
		}

		//// ... CRLF = CR = 0x0d + LF = 0x0a
		//netData, err := bufio.NewReader(conn).ReadBytes(0x0a)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//
		//t := strings.TrimSpace(string(netData))
		//
		//if strings.Contains(t, EHLO) || strings.Contains(t, "ehlo") {
		//	conn.Write([]byte(R_EHLO))
		//}
		//
		//if strings.Contains(t, HELO) || strings.Contains(t, "helo") {
		//	conn.Write([]byte(R_HELO_DOMAIN))
		//}
		//
		//if strings.Contains(t, MAIL_FROM) || strings.Contains(t, "mail FROM:") {
		//	fmt.Println("FROM MAIL: ", parsCommand(t))
		//	conn.Write([]byte(R_MAIL_FROM))
		//}
		//
		//if strings.Contains(t, RCPT) || strings.Contains(t, "rcpt TO:") {
		//	fmt.Println("RCPT TO: ", parsCommand(t))
		//	conn.Write([]byte(R_RCPT))
		//}
		//
		//if t == DATA || t == "data" {
		//	conn.Write([]byte(R_DATA))
		//}
		//
		//if bytes.Contains(netData, []byte{0x2e, 0x0d, 0x0a}) {
		//	conn.Write([]byte(R_CRLF_POINT_CRLF))
		//	fmt.Printf("%x", t)
		//}
		//
		//if t == QUIT || t == "quit" {
		//	conn.Write([]byte(R_QUIT))
		//	conn.Close()
		//	w.Done()
		//}
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
