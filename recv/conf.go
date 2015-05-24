package main

import (
	"crypto/tls"
)

type conf struct {
	lAddr      string
	srvName    string
	extensions []string
	tls        *tls.Config
}

var certfile, keyfile = "/home/dan/code/go/src/znvy/certs/cert.pem",
	"/home/dan/code/go/src/znvy/certs/key.pem"

var cert tls.Certificate
var config conf

func init() {
	var err error
	cert, err = tls.LoadX509KeyPair(certfile, keyfile)
	xk(err, "Loading certificate and key files: %s %s", certfile, keyfile)
	config = conf{
		lAddr:      "127.0.0.1:1234",
		srvName:    "bronco.here",
		extensions: []string{"STARTTLS"},
		tls: &tls.Config{
			Certificates: []tls.Certificate{cert},
			InsecureSkipVerify: true,
			ServerName: "bronco.here",
		},
	}
}
