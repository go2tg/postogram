package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "2525"

	GREETINGS = "mail.go-tg.test SMTP is glad to see you!\n"
	HELO      = "HELO"
	R_DOMAIN  = "250 domain name should be qualified\n"
	MAIL_FROM = "MAIL FROM:"
	R_        = "250 someusername@somecompany.ru sender accepted\n"
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

		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		t := strings.TrimSpace(netData)

		if t == HELO {
			conn.Write([]byte(R_DOMAIN))
		}

		if strings.Contains(t, MAIL_FROM) {
			strings.Trim(t, ":")
			fmt.Println(t)
			re := regexp.MustCompile("\\<(.*?)\\>")
			match := re.FindStringSubmatch(t)
			fmt.Println(match[1])
		}

		//switch strings.TrimSpace(netData) {
		//case HELO:
		//	conn.Write([]byte(R_DOMAIN))
		//case (strings.Contains(MAIL_FROM)) :
		//	conn.Write([]byte(R_))
		//}
	}

	conn.Close()
}

//func HandSend(conn net.Conn) {
//
//	if GREETINGS_COUNT == 1 {
//		conn.Write([]byte(GREETINGS))
//		GREETINGS_COUNT++
//	}
//
//	netData, err := bufio.NewReader(conn).ReadString('\n')
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	switch strings.TrimSpace(netData) {
//	case HELO:
//		conn.Write([]byte(R_DOMAIN))
//	}
//
//
//	//temp := strings.TrimSpace(netData)
//	//if temp == HELO {
//	//	conn.Write([]byte(R_DOMAIN))
//	//	return
//	//}
//}
