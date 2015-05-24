package main

import (
	rex "regexp"
)

var (
	reEhlo = rex.MustCompile(`EHLO [^@]+`)
	reStartTLS = rex.MustCompile(`STARTTLS`)
	reFrom = rex.MustCompile(`MAIL FROM:(.*)`)
	reTo = rex.MustCompile(`RCPT TO:(.*)`)
)

// expect looks for a match in the coming line of input, returning
// output of re.FindSubmatch if successful. Errors are strictly from
// the read operation: a no-match does not yield an error.
func (h handler) expect(re *rex.Regexp) (matches [][]byte, err error) {
	l, err := h.ReadLineBytes()
	if err != nil {
		return
	}
	matches = re.FindSubmatch(l)
	return
}

// until wraps expect by looping it and informing the client if the 
// input command is not what was expected. Any error besides a no-match
// is bubbled up. 
// TODO some other error handling here is a good idea, as well 
// as adding a timeout.
func (h handler) until(re *rex.Regexp) (matches [][]byte, err error) {
	for {
		matches, err =  h.expect(re)
		if err != nil || matches != nil {
			return
		}
		h.reply(503) //kind of a filler response
	}
}
