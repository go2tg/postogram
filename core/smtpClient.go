package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("localhost:2525")
	if err != nil {
		fmt.Println(err)
	}

	// Set the sender and recipient first
	if err := c.Mail("sender@example.org"); err != nil {
		fmt.Println(err)
	}
	if err := c.Rcpt("recipient@example.net"); err != nil {
		fmt.Println(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		fmt.Println(err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		fmt.Println(err)
	}
	err = wc.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
