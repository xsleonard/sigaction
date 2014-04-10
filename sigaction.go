// Package sigaction provides a method for binding a function to any number
// of os.Signals. It also provides some convenience methods for common bindings.
// All of the os.Signals are actually in pacakge syscall, but package os
// re-exports syscall.SIGINT and syscall.SIGKILL as os.Interrupt and os.Kill.
package sigaction

import (
	"log"
	"os"
	"os/signal"
)

// Calls a function when signals are caught. You will want to run this as a
// goroutine except when you know you don't.
func SigAction(f func(), signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	<-c
	f()
}

// Panics on SIGINT (CTRL+C)
func InterruptPanic() {
	SigAction(func() { log.Panic() }, os.Interrupt)
}
