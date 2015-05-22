package main

import (
	"net"
	txt "net/textproto"
	"time"
	"znvy/lg"
)

var ck, xk, log = lg.Ck, lg.Xk, lg.Std

func main() {
	l, err := net.Listen("tcp", config.lAddr)
	xk(err, "Listening on %s", config.lAddr)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if ck(err, "Accepting  connection") {
			continue
		}

		log.Println("Connection accepted: ", conn.RemoteAddr())

		conn.SetDeadline(time.Now().Add(20 * time.Second))

		go func(conn net.Conn) {
			h := handler{
				txt.NewConn(conn),
				conn,
			}
			h.handle()
		}(conn)
	}
}
