package lg

import (
	"fmt"
	"log"
	"os"
)

var Ck, Xk = Std.Ck, Std.Xk
var Std = Logger{ log.New(os.Stdout, "", 0) }

// Logger is log.Logger with a couple conveniences
type Logger struct { *log.Logger }

// Ck checks an error, prints any message that has been passed in, 
// and returns err != nil
func (l Logger) Ck(err error, msgs ...interface{}) bool {
	if err == nil {
		if len(msgs) != 0 {
			l.Printf("%s... good.\n", 
				fmt.Sprintf(msgs[0].(string), msgs[1:]...))
		}
		return false
	} else {
		if len(msgs) != 0 {
			l.Printf("%s... FAIL. '%v'\n", 
				fmt.Sprintf(msgs[0].(string), msgs[1:]...), 
				err)
		}
		return true
	}
}

// Xk performs the same operations as Ck, but if err != nil,
// it initiates a panic
func (l Logger) Xk(err error, msgs ...interface{}) {
	if l.Ck(err, msgs...) {
		panic(err)
	}
}
