package main

import (
	"crypto/tls"
	"net"
	txt "net/textproto"
	"net/mail"
)

// handler negotiates SMTP receipt
type handler struct {
	*txt.Conn
	netc net.Conn //hang on to the unwrapped connection for some behaviour
}

// handle() wraps the entire negotiation from beginnning to end.
// If any sub-neogtiation returns an error, handle closes the connection
// and returns.
func (h handler) handle() {
	var (
		match [][]byte
		err error
	)
	defer h.Close()

	raddr := h.netc.RemoteAddr()

	if ck(h.PrintfLine("220 %s ZNVY", config.srvName), "<%s> Greeting",raddr) {
		h.reply(666)
		return
	}

	// EHLO
	match, err = h.until(reEhlo)
	if ck(err, "<%s> waiting for EHLO", raddr) {
		h.reply(666)
		return
	}
	//ehloName := match[1]

	h.PrintfLine("%d-%s", 250, config.srvName)
	for i, s := range config.extensions {
		if i == len(config.extensions)-1 {
			h.msg(250, s)
			break
		}
		h.PrintfLine("%d-%s", 250, s)
	}

	// STARTTLS
	match, err = h.until(reStartTLS)
	if ck(err, "<%s> waiting for StartTLS", raddr) {
		h.reply(666)
		return
	}
	h.msg(220, "Go ahead")

	// actually negotiate it
	tlsc := tls.Server(h.netc, config.tls)
	if ck(tlsc.Handshake(), "<%s> TLS handshake", raddr) {
		h.msg(403, "TLS handshake failed")
		return
	}
	h.Conn = txt.NewConn(tlsc)

	// another EHLO
	match, err = h.until(reEhlo)
	if ck(err, "<%s> EHLO after StartTLS", raddr) {
		h.reply(666)
		return
	}

	// MAIL FROM
	match, err = h.until(reFrom)
	if ck(err, "<%s> MAIL FROM", raddr) {
		h.reply(666)
		return
	}
	user, domain := match[1], match[2]
	log.Printf("<%s> MAIL FROM %s@%s\n", raddr, user, domain)

	//RCPT TO - the first time
	match, err = h.until(reTo)
	if ck(err, "<%s> RCPT TO", raddr) {
		h.reply(666)
		return
	}



	//optional additional RCPT's
	for {
		mat
	}
}

// send a numeric reply + msg
func (h handler) msg(code int, msg string) {
	h.PrintfLine("%d %s", code, msg)
}

// generic messages from codes
func (h handler) reply(code int) {
	var msg string
	switch code {
	case 250:
		msg = "OK"
	case 451:
		msg = "Timeout"
	case 503:
		msg = "Bad sequence"
	case 530:
		msg = "Must issue a STARTTLS command first"
	case 666:
		msg = "Unspecified error"
	}
	h.msg(code, msg)
}
