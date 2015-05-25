package main

import (
	"fmt"
)

type smtpResponse int

const (
	OK            smtpResponse = 250
	LOCAL_ERR     smtpResponse = 451
	SYNTAX_CMD    smtpResponse = 500
	SYNTAX_PARAM  smtpResponse = 501
	SEQUENCE      smtpResponse = 503
	MUST_STARTTLS smtpResponse = 530
)

func (r smtpResponse) String() string {
	return fmt.Sprintf("%d %s", r, repsonseStrings[r])
}

func (r smtpResponse) Error() string {
	return r.String()
}

var repsonseStrings = map[smtpResponse]string{
	250: "OK",
	451: "Action aborted: local error in processing", // catch-all
	500: "Bad syntax: command unrecognized",
	501: "Bad syntax in parameter",
	503: "Bad sequence of commands",
	530: "Must issue a STARTTLS command first",
}
