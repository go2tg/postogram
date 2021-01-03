package main

import (
	"bytes"
	"errors"
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
	QUIT_PY           = "quit"
	R_QUIT            = "221 \n"
)

//var CRLFPointCRLF = []byte{0x2e, 0x0d, 0x0a}

type MailSession struct {
	session  string
	ehlo     bool
	helo     bool
	mailFrom string
	rcptTo   string
	data     bool
	dataBody []byte
	endData  bool
	quit     bool
}

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
		go handleConnection(conn) // при QUIT завершается рутина
	}
}

func handleConnection(conn net.Conn) {

	var mail MailSession
	// Отправка приветствия
	conn.Write([]byte(GREETINGS))
	buf := make([]byte, 128)

	//for {
	fmt.Println(mail.quit)

L1:

	if mail.quit == false {
		// Чтение из сети
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error read TCP   : ", err)
		}
	}

	//t := strings.TrimSpace(string(buf))

	b, msg, err := CRLF(buf)
	if err != nil {
		fmt.Println(err)
	} else {
		if b {
			fmt.Printf("%s", msg)
			fmt.Println("")

			// EHLO
			if bytes.Equal([]byte(EHLO), msg[0:4]) || bytes.Equal([]byte(EHLO_PY), msg[0:4]) {
				fmt.Println("EHLO - OK")
				conn.Write([]byte(R_EHLO))
				mail.ehlo = true
				fmt.Println(buf)
				goto L1
			}

			// HELO
			if bytes.Equal([]byte(HELO), msg[0:4]) || bytes.Equal([]byte(HELO_PY), msg[0:4]) {
				fmt.Println("HELO - OK")
				conn.Write([]byte(R_HELO_DOMAIN))
				mail.helo = true
				fmt.Println(buf)
				goto L1
			}

			// MAIL FROM:
			if bytes.Equal([]byte(MAIL_FROM), msg[0:10]) || bytes.Equal([]byte(MAIL_FROM_PY), msg[0:10]) {
				t := strings.TrimSpace(string(msg))
				fmt.Println("MAIL FROM: - OK ", parsCommand(t))
				conn.Write([]byte(R_MAIL_FROM))
				mail.mailFrom = parsCommand(t)
				fmt.Println(buf)
				goto L1
			}

			// RCPT TO:
			if bytes.Equal([]byte(RCPT), msg[0:8]) || bytes.Equal([]byte(RCPT_PY), msg[0:8]) {
				t := strings.TrimSpace(string(msg))
				fmt.Println("RCPT TO: - OK ", parsCommand(t))
				conn.Write([]byte(R_RCPT))
				mail.rcptTo = parsCommand(t)
				fmt.Println(buf)
				goto L1
			}
			//

			// DATA
			if bytes.Equal([]byte(DATA), msg[0:4]) || bytes.Equal([]byte(DATA_PY), msg[0:4]) {
				fmt.Println("DATA - OK")
				conn.Write([]byte(R_DATA))
				mail.data = true
				fmt.Println(buf)
				goto L1
			}

			in, mes, err := CRLFPointCRLF(buf)
			if err != nil {
				fmt.Println(err)
			} else {
				if in {
					fmt.Println(mes)
					mail.endData = true
					fmt.Println(buf, "--")
					conn.Write([]byte(R_CRLF_POINT_CRLF))
					//buf = buf[:0]
					//goto L1

					_, err = conn.Read(buf)
					if err != nil {
						fmt.Println("Error read TCP   : ", err)
					}

					if bytes.Equal([]byte(QUIT), buf[0:4]) || bytes.Equal([]byte(QUIT_PY), buf[0:4]) {
						conn.Write([]byte(R_QUIT))
						//fmt.Println(buf)

						mail.quit = true
						conn.Close()
						fmt.Printf("%v", mail)
					}

				}
			}

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

func CRLF(message []byte) (b bool, mes []byte, err error) {

	if cap(message) < 2 {
		return false, nil, errors.New("message len < 2")
	} else {
		for i, v := range message {
			if v == 0x0d && message[i+1] == 0x0A {
				fmt.Println("CRLF")
				fmt.Println("Position", i, i+1)
				mes = message[0:i]
				return true, mes, nil
			}
		}
	}
	return false, nil, nil
}

func CRLFPointCRLF(message []byte) (b bool, mes []byte, err error) {

	if cap(message) < 5 {
		return false, nil, errors.New("message len < 2")
	} else {
		for i, v := range message {
			if v == 0x0d && message[i+1] == 0x0a && message[i+2] == 0x2e && message[i+3] == 0x0d && message[i+4] == 0x0a {
				fmt.Println("CRLF . CRLF")
				fmt.Println("Position", i, "-", i+4)
				mes = message[0:i]
				return true, mes, nil
			}
		}
	}
	return false, nil, nil
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
//}
