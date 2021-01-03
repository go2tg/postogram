package main

import (
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go XMail(&wg)
	}
	stop := time.Now()
	wg.Wait()
	fmt.Printf("Elapse time: %v ", stop.Sub(start))

}

func XMail(wg *sync.WaitGroup) {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("localhost:2525")
	if err != nil {
		fmt.Println("TCP connect error", err)
	}

	// Set the sender and recipient first
	if err = c.Mail("sender@example.org"); err != nil {
		fmt.Println("Sender set error", err)
	}
	if err = c.Rcpt("recipient@example.net"); err != nil {
		fmt.Println("RCPT set error", err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		fmt.Println(err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body. 1000, 1000 23/ end")
	if err != nil {
		fmt.Println("Set BODY e-mail", err)
	}
	err = wc.Close()
	if err != nil {
		fmt.Println("Close  error", err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal("QUIT error", err)
	}
	wg.Done()

}
