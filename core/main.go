package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"
)

const (
	HOST           = "localhost"
	PORT           = "2525"
	GREETINGS      = "220 localhost.local\n"
	EHLO           = "EHLO"
	EhloPy         = "ehlo"
	REhlo          = "500 \n"
	HELO           = "HELO"
	HeloPy         = "helo"
	RHeloDomain    = "250 \n"
	MailFrom       = "MAIL FROM:"
	MailFromPy     = "mail FROM:"
	RMailFrom      = "250 \n"
	RCPT           = "RCPT TO:"
	RcptPy         = "rcpt TO:"
	RRcpt          = "250 \n"
	DATA           = "DATA"
	DataPy         = "data"
	RData          = "354 \n"
	RCrlfPointCrlf = "250 \n"
	QUIT           = "QUIT"
	QuitPy         = "quit"
	RQuit          = "221 \n"
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

	cpuFile, err := os.Create("/tmp/cpuProfile.out")
	if err != nil {
		fmt.Println(err)
		return
	}
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Port error", err.Error())
		return
	}

	defer l.Close()

	//for x := 0 ; x <1000 ; x++{
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Port error", err.Error())
		l.Close()
		return
		//	}
		go handleConnection(conn) // при QUIT завершается рутина
	}
}

func handleConnection(conn net.Conn) {

	var mail MailSession
	// Отправка приветствия
	_, _ = conn.Write([]byte(GREETINGS))
	buf := make([]byte, 1024)

	//for {
	//fmt.Println(mail.quit)
	mail.session = GenUUIDv4()

L1:

	if mail.quit == false {
		// Чтение из сети
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error read TCP   : ", err)
			conn.Close()
		}
	}

	// Фильтруем команды по разделителю CRLF
	b, msg, err := CRLF(buf)
	if err != nil {
		fmt.Println(err)
	} else {
		if b {
			//fmt.Printf("%s", msg)
			//fmt.Println("")

			// EHLO
			if bytes.Equal([]byte(EHLO), msg[0:4]) || bytes.Equal([]byte(EhloPy), msg[0:4]) {
				//fmt.Println("EHLO - OK")
				_, _ = conn.Write([]byte(REhlo))
				mail.ehlo = true
				goto L1
			}

			// HELO
			if bytes.Equal([]byte(HELO), msg[0:4]) || bytes.Equal([]byte(HeloPy), msg[0:4]) {
				//fmt.Println("HELO - OK")
				_, _ = conn.Write([]byte(RHeloDomain))
				mail.helo = true
				goto L1
			}

			// MAIL FROM:
			if bytes.Equal([]byte(MailFrom), msg[0:10]) || bytes.Equal([]byte(MailFromPy), msg[0:10]) {
				t := strings.TrimSpace(string(msg))
				//fmt.Println("MAIL FROM: - OK ", parsCommand(t))
				_, _ = conn.Write([]byte(RMailFrom))
				mail.mailFrom = ParsCommand(t)
				goto L1
			}

			// RCPT TO:
			if bytes.Equal([]byte(RCPT), msg[0:8]) || bytes.Equal([]byte(RcptPy), msg[0:8]) {
				t := strings.TrimSpace(string(msg))
				//fmt.Println("RCPT TO: - OK ", parsCommand(t))
				_, _ = conn.Write([]byte(RRcpt))
				mail.rcptTo = ParsCommand(t)
				goto L1
			}
			//

			// DATA
			if bytes.Equal([]byte(DATA), msg[0:4]) || bytes.Equal([]byte(DataPy), msg[0:4]) {
				//fmt.Println("DATA - OK")
				_, _ = conn.Write([]byte(RData))
				mail.data = true
				goto L1
			}

			//fmt.Println(buf)
			//fmt.Println(string(buf))

			// Отсекаем окончание блока сообщение по CRLF.CRLF
			in, mes, err := CRLFPointCRLF(buf)
			if err != nil {
				fmt.Println(err)
			} else {
				if in {
					mail.dataBody = mes
					mail.endData = true
					_, _ = conn.Write([]byte(RCrlfPointCrlf))
					//buf = buf[:0]
					//goto L1

					_, err = conn.Read(buf)
					if err != nil {
						fmt.Println("Error read TCP   : ", err)
						conn.Close()
					}

					//fmt.Println(buf)
					if bytes.Equal([]byte(QUIT), buf[0:4]) || bytes.Equal([]byte(QuitPy), buf[0:4]) {
						_, _ = conn.Write([]byte(RQuit))
						mail.quit = true
						conn.Close()
						//fmt.Printf("%v", mail)
					}
				}
			}

		}
	}

	//fmt.Printf("\n")
	//fmt.Printf("%v", mail)
	//fmt.Printf("\n\t\n")
	//conn.Close()
}

func ParsCommand(t string) string {
	strings.Trim(t, ":")
	re := regexp.MustCompile("\\<(.*?)\\>")
	match := re.FindStringSubmatch(t)
	if match != nil {
		return match[1]
	} else {
		return ""
	}
}

// CRLF - поиск в массиве байт последовательности  0x0d , 0x0a
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

// CRLFPointCRLF - поиск в массиве байт последовательности 0x0d , 0x0a , 0x2e , 0x0d , 0x0a
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
	_, _ = rand.Read(u)
	//Set the version to 4
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	ss := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return ss
}
