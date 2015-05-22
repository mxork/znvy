package main

import (
	"fmt"
	"net/mail"
	"os"
	"znvy/lg"
)

var ck = lg.Ck

func main() {
	fname := os.Args[1]
	f, err := os.Open(fname)
	ck(err, "Opening test message")

	msg, err :=	mail.ReadMessage(f)
	ck(err, "Reading message")

	hd, _ := msg.Header, msg.Body

	pretty(hd)

}

func pretty(hd mail.Header) {
	for k, v := range hd {
		fmt.Println(k)
		for _, vv := range v {
			fmt.Println("\t" + vv)
		}
	}
}

