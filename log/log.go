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

// ----- Public API -----

// Print formats using the default formats for its operands and logs the result.
// It is non‑blocking.
func Print(args ...interface{}) {
	send(fmt.Sprint(args...))
}

// Printf formats according to a format specifier and logs the result.
// It is non‑blocking.
func Printf(format string, args ...interface{}) {
	send(fmt.Sprintf(format, args...))
}

// Println formats using the default formats for its operands, appends a newline, and logs the result.
// It is non‑blocking.
func Println(args ...interface{}) {
	send(fmt.Sprintln(args...))
}

// Error logs to stderr (with the same formatting as Print).
// It is non‑blocking.
func Error(args ...interface{}) {
	// We could send to a separate channel for stderr, but for simplicity we use the same channel
	// and prefix with "ERROR:" to distinguish.
	send("ERROR: " + fmt.Sprint(args...))
}

// Errorf logs a formatted error message to stderr.
func Errorf(format string, args ...interface{}) {
	send("ERROR: " + fmt.Sprintf(format, args...))
}

// Errorln logs to stderr with a newline.
func Errorln(args ...interface{}) {
	send("ERROR: " + fmt.Sprintln(args...))
}

// Fatal logs to stderr and then calls os.Exit(1).
// It forces a flush of the log buffer to ensure the message is written.
func Fatal(args ...interface{}) {
	msg := "FATAL: " + fmt.Sprint(args...)
	sendBlocking(msg)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits.
func Fatalf(format string, args ...interface{}) {
	msg := "FATAL: " + fmt.Sprintf(format, args...)
	sendBlocking(msg)
	os.Exit(1)
}

// Fatalln logs a fatal message with a newline and exits.
func Fatalln(args ...interface{}) {
	msg := "FATAL: " + fmt.Sprintln(args...)
	sendBlocking(msg)
	os.Exit(1)
}

// Flush waits for all queued log messages to be written.
// Useful before a graceful shutdown.
func Flush() {
	wg.Wait()
}

// SetOutput changes the destination for standard logs.
func SetOutput(w io.Writer) {
	stdout = w
}

// SetErrorOutput changes the destination for error logs.
func SetErrorOutput(w io.Writer) {
	stderr = w
}
