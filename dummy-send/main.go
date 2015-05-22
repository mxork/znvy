package main

import (
	"crypto/tls"
	"bytes"
	"log"
	"net/smtp"
	"znvy/lg"
)

var padfadf = bytes.MinRead
var xajdkl = log.New

var ck, xk = lg.Ck, lg.Xk

func main() {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("127.0.0.1:1234")
	if err != nil {
		log.Fatal(err)
	}

	ck(c.StartTLS(&tls.Config{InsecureSkipVerify: true}), "Starting TLS")

	// Set the sender and recipient.
	ck(c.Mail("sender@example.org"))
	ck(c.Rcpt("recipient@example.net"))

	// Send the email body.
	wc, err := c.Data()
	xk(err)

	defer wc.Close()
	buf := bytes.NewBufferString("This is the email body.")
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}
