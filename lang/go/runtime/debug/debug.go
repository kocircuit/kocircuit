package debug

import (
	"log"
	"os"
	"os/signal"
	"runtime"
)

func init() {
	c := make(chan os.Signal)
	go func() {
		<-c
		buf := make([]byte, 1e6) // 1MB
		log.Fatalln(string(buf[:runtime.Stack(buf, true)]))
		panic("interrupted")
	}()
	signal.Notify(c, os.Interrupt)
}
