package g

import (
	"log"
	"runtime"
)

const (
	VERSION = "1.0.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
