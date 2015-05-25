package main

type command int
const (
	EHLO command = iota
	STARTTLS
	AUTH
	MAIL_FROM
	RCPT_TO
	DATA
	QUIT
	RSET
)

func (cmd command) String() string {
	return cmdStrings[cmd]
}

var cmdStrings = [...]string{
	"EHLO",
	"STARTTLS",
	"AUTH",
	"MAIL FROM:",
	"RCPT TO:",
	"DATA",
	"QUIT",
	"RSET",
}

func match(prefix string, line []byte) bool {
	return string(line[len(prefix)]) == prefix
}

func cmdLookup(line []byte) command {
	for i, prefix := range cmdStrings{
		if match(prefix, line) {
			return command(i)
		}
	}

	return -1
}
