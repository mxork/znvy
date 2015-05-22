package main

import (
	"crypto/tls"
	"net"
	txt "net/textproto"
	re "regexp"
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
	defer h.Close()
	raddr := h.netc.RemoteAddr()

	if ck(h.greet(), "<%s> Greeting",raddr) {
		h.reply(666)
		return
	}

	if ck(h.expectHello(), "<%s> EHLO", raddr) {
		h.reply(666)
		return
	}

	if ck(h.expectStartTLS(), 
			"<%s> STARTTLS",
			raddr) {
		h.reply(666)
		return
	}

	// actually negotiate it
	tlsc := tls.Server(h.netc, config.tls)
	tlsc.Handshake()
	h.Conn = txt.NewConn(tlsc)

	if ck(h.expectHello(),
			"<%s> EHLO after STARTTLS",
			raddr) {
		h.reply(666)
		return
	}
	addr, user, domain, err := h.expectFrom()
	ck(err, "FROM %s %s %s", addr, user, domain)
}

// Normal course
func (h handler) greet() error {
	return h.PrintfLine("220 %s ZNVY", config.srvName)
}

// Expects
var rehelo = re.MustCompile(`EHLO [^@]+`)

func (h handler) expectHello() error {
	for {
		l, err := h.next()
		if err != nil {
			return err // pop the err up the line
		}

		res := rehelo.FindSubmatch(l)
		if res == nil {
			h.reply(666)
			continue
		}

		// success...
		h.PrintfLine("%d-%s", 250, config.srvName)
		for i, s := range config.extensions {
			if i == len(config.extensions)-1 {
				h.PrintfLine("%d %s", 250, s)
				break
			}
			h.PrintfLine("%d-%s", 250, s)
		}

		return nil
	}
}

var restarttls = re.MustCompile(`STARTTLS`)

// currently junk
func (h handler) expectStartTLS() error {
	for {
		l, err := h.next()
		if err != nil {
			return err
		}

		l = txt.TrimBytes(l)

		ok := restarttls.Match(l)
		if !ok {
			log.Print("Failed to match for startTLS")
			h.reply(530)
			continue
		}

		h.msg(220, "Ready to start TLS")

		return nil
	}
}

var refrom = re.MustCompile(`MAIL FROM:\s*([^@]+)@([^@]+)`)

func (h handler) expectFrom() (addr, user, domain string, err error) {
	for {
		l, err := h.next()
		if ck(err) {
			return "", "", "", err
		}

		loc := refrom.FindSubmatchIndex(l)
		if loc == nil {
			h.reply(666)
			continue
		}

		h.reply(250)

		// user and domain indices
		u0, u1, d0, d1 := loc[2], loc[3], loc[4], loc[5]

		// return all, and user/dom seperate
		return string(l[u0:d1]), string(l[u0:u1]), string(l[d0:d1]), nil
	}
}

func (h handler) expectFor() {

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

// next() fetches the next line of input
func (h handler) next() (line []byte, err error) {
	line, err = h.ReadLineBytes()
	// TODO wrap some trivial errors and keep going
	return
}

// not in use
func (h handler) parseCmd() {
	ln, err := h.next()
	ck(err, "Next line of input")
	ln = txt.TrimBytes(ln)

	// Essentially strings.SplitOnWs

	argv := [][]byte{}

	first, inword := 0, false
	for i, v := range string(ln) {
		isWS := v == ' ' || v == '\t'

		if inword && isWS {
			inword = false
			argv = append(argv, ln[first:i])
			continue
		}

		if !inword && !isWS {
			inword = true
			first = i
			continue
		}

	}

	// flush the last word, if there
	if inword {
		argv = append(argv, ln[first:])
	}

	// checking against supported commands
	cmd := argv[0]

	switch string(cmd) {
	case "STARTTLS":
	default:
		//unsupported
	}
}
