package main

import (
	"bytes"
	"crypto/rand"
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
	//fmt.Println(mail.quit)
	mail.session = GenUUIDv4()

L1:

	if mail.quit == false {
		// Чтение из сети
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error read TCP   : ", err)
		}
	}

	b, msg, err := CRLF(buf)
	if err != nil {
		fmt.Println(err)
	} else {
		if b {
			//fmt.Printf("%s", msg)
			//fmt.Println("")

			// EHLO
			if bytes.Equal([]byte(EHLO), msg[0:4]) || bytes.Equal([]byte(EHLO_PY), msg[0:4]) {
				fmt.Println("EHLO - OK")
				conn.Write([]byte(R_EHLO))
				mail.ehlo = true
				goto L1
			}

			// HELO
			if bytes.Equal([]byte(HELO), msg[0:4]) || bytes.Equal([]byte(HELO_PY), msg[0:4]) {
				fmt.Println("HELO - OK")
				conn.Write([]byte(R_HELO_DOMAIN))
				mail.helo = true
				goto L1
			}

			// MAIL FROM:
			if bytes.Equal([]byte(MAIL_FROM), msg[0:10]) || bytes.Equal([]byte(MAIL_FROM_PY), msg[0:10]) {
				t := strings.TrimSpace(string(msg))
				fmt.Println("MAIL FROM: - OK ", parsCommand(t))
				conn.Write([]byte(R_MAIL_FROM))
				mail.mailFrom = parsCommand(t)
				goto L1
			}

			// RCPT TO:
			if bytes.Equal([]byte(RCPT), msg[0:8]) || bytes.Equal([]byte(RCPT_PY), msg[0:8]) {
				t := strings.TrimSpace(string(msg))
				fmt.Println("RCPT TO: - OK ", parsCommand(t))
				conn.Write([]byte(R_RCPT))
				mail.rcptTo = parsCommand(t)
				goto L1
			}
			//

			// DATA
			if bytes.Equal([]byte(DATA), msg[0:4]) || bytes.Equal([]byte(DATA_PY), msg[0:4]) {
				fmt.Println("DATA - OK")
				conn.Write([]byte(R_DATA))
				mail.data = true
				goto L1
			}

			in, mes, err := CRLFPointCRLF(buf)
			if err != nil {
				fmt.Println(err)
			} else {
				if in {
					fmt.Println("Message body - OK")
					fmt.Println("message :", mes)
					mail.dataBody = mes
					mail.endData = true
					conn.Write([]byte(R_CRLF_POINT_CRLF))
					//buf = buf[:0]
					//goto L1

					_, err = conn.Read(buf)
					if err != nil {
						fmt.Println("Error read TCP   : ", err)
					}

					if bytes.Equal([]byte(QUIT), buf[0:4]) || bytes.Equal([]byte(QUIT_PY), buf[0:4]) {
						conn.Write([]byte(R_QUIT))
						mail.quit = true
						conn.Close()
						//fmt.Printf("%v", mail)
					}

				}
			}

		}
	}

	fmt.Printf("\n")
	fmt.Printf("%v", mail)
	fmt.Printf("\n\t\n")
	//conn.Close()
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
				mes = message[0:i]
				return true, mes, nil
			}
		}
	}
	return false, nil, nil
}

// GenUUIDv4 - генерирует UUID v4
func GenUUIDv4() string {
	u := make([]byte, 16)
	rand.Read(u)
	//Set the version to 4
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	ss := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return ss
}
