package sigaction

import (
	"os"
	"testing"
	"time"
)

func raise(t *testing.T, sig os.Signal) {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}
	process.Signal(sig)
}

func TestSigAction(t *testing.T) {
	signals := []os.Signal{os.Interrupt, os.Kill}
	sig := os.Interrupt
	// Can't test with os.Kill since its not subduable and will make tests fail
	triggered := make(chan bool)
	go func() {
		SigAction(func() {
			close(triggered)
		}, signals...) // pass in all signals to verify multiargs
	}()

	time.Sleep(time.Millisecond * 10)
	raise(t, sig)
	time.Sleep(time.Millisecond * 10)
	tick := time.Tick(time.Second)
	select {
	case <-triggered:
		break
	case <-tick:
		t.Fatalf("signal %v wasn't triggered in time", sig)
		break
	}
}

func TestInterruptPanic(t *testing.T) {
	go func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("Expected panic")
			}
		}()
		InterruptPanic()
	}()
	time.Sleep(time.Millisecond * 10)
	raise(t, os.Interrupt)
}
