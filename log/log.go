// Package dukelog helps in async logging
package dukelog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const defaultBufferSize = 10000

var (
	logCh = make(chan string, defaultBufferSize)

	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr

	wg sync.WaitGroup
)

func init() {
	go func() {
		for msg := range logCh {
			fmt.Fprintln(stdout, msg)
			wg.Done() // signal that this message has been written
		}
	}()
}

// internal helper to send a message asynchronously.
// It increments the WaitGroup before sending, so we can wait for flush.
// The send is non‑blocking; if the buffer is full, the message is dropped.
func send(msg string) {
	wg.Add(1)
	select {
	case logCh <- msg:
		// message queued
	default:
		// buffer full – drop the message
		wg.Done() // we incremented but didn't send, so we must decrement
	}
}

func sendBlocking(msg string) {
	wg.Add(1)
	logCh <- msg // blocks until there is room in the buffer
	wg.Wait()    // wait for the writer to process it
}

func Print(args ...interface{}) {
	send(fmt.Sprint(args...))
}

func Printf(format string, args ...interface{}) {
	send(fmt.Sprintf(format, args...))
}

func Println(args ...interface{}) {
	send(fmt.Sprintln(args...))
}

func Error(args ...interface{}) {
	send("ERROR: " + fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	send("ERROR: " + fmt.Sprintf(format, args...))
}

func Errorln(args ...interface{}) {
	send("ERROR: " + fmt.Sprintln(args...))
}

func Fatal(args ...interface{}) {
	msg := "FATAL: " + fmt.Sprint(args...)
	sendBlocking(msg)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	msg := "FATAL: " + fmt.Sprintf(format, args...)
	sendBlocking(msg)
	os.Exit(1)
}

func Fatalln(args ...interface{}) {
	msg := "FATAL: " + fmt.Sprintln(args...)
	sendBlocking(msg)
	os.Exit(1)
}

func Flush() {
	wg.Wait()
}

func SetOutput(w io.Writer) {
	stdout = w
}

func SetErrorOutput(w io.Writer) {
	stderr = w
}
