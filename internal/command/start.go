package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
)

// Start launches a command and streams stdout and stderr to channels.
// All the channels returned are ready only and won't be closed
// if the command fails later.
func (c *Cmder) Start(cmd *exec.Cmd) (
	stdoutLines, stderrLines <-chan string,
	waitError <-chan error, startErr error,
) {
	return start(cmd)
}

func start(cmd execCmd) (stdoutLines, stderrLines <-chan string,
	waitError <-chan error, startErr error,
) {
	stop := make(chan struct{})
	stdoutReady := make(chan struct{})
	stdoutLinesCh := make(chan string)
	stdoutDone := make(chan struct{})
	stderrReady := make(chan struct{})
	stderrLinesCh := make(chan string)
	stderrDone := make(chan struct{})

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	go streamToChannel(stdoutReady, stop, stdoutDone, stdout, stdoutLinesCh)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		close(stop)
		<-stdoutDone
		return nil, nil, nil, err
	}
	go streamToChannel(stderrReady, stop, stderrDone, stderr, stderrLinesCh)

	err = cmd.Start()
	if err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		close(stop)
		<-stdoutDone
		<-stderrDone
		return nil, nil, nil, err
	}

	waitErrorCh := make(chan error)
	go func() {
		err := cmd.Wait()
		_ = stdout.Close()
		_ = stderr.Close()
		close(stop)
		<-stdoutDone
		<-stderrDone
		waitErrorCh <- err
	}()

	return stdoutLinesCh, stderrLinesCh, waitErrorCh, nil
}

func streamToChannel(ready chan<- struct{},
	stop <-chan struct{}, done chan<- struct{},
	stream io.Reader, lines chan<- string,
) {
	defer close(done)
	close(ready)
	scanner := bufio.NewScanner(stream)
	lineBuffer := make([]byte, bufio.MaxScanTokenSize) // 64KB
	const maxCapacity = 20 * 1024 * 1024               // 20MB
	scanner.Buffer(lineBuffer, maxCapacity)

	for scanner.Scan() {
		// scanner is closed if the context is canceled
		// or if the command failed starting because the
		// stream is closed (io.EOF error).
		lines <- scanner.Text()
	}
	err := scanner.Err()
	if err == nil || errors.Is(err, os.ErrClosed) {
		return
	}

	// ignore the error if it is stopped.
	select {
	case <-stop:
		return
	default:
		lines <- "stream error: " + err.Error()
	}
}
