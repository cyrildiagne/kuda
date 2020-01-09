package main

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os/exec"
)

// RunSkaffold launches 'run' on skaffold and streams logs to w.
func RunSkaffold(tempDir string, skaffoldFile string, w io.Writer) error {
	// Run Skaffold Deploy.
	args := []string{"run", "-f", skaffoldFile}
	cmd := exec.Command("skaffold", args...)
	cmdout, _ := cmd.StdoutPipe()
	cmderr, _ := cmd.StderrPipe()
	cmd.Dir = tempDir

	if err := cmd.Start(); err != nil {
		return err
	}

	// Make sure the client accepts server side events.
	conn, ok := w.(http.Flusher)
	if !ok {
		return errors.New("clients must support server side events")
	}

	errc := make(chan error)

	go func() {
		err := copyAndFlush(w, cmdout, conn)
		errc <- err
	}()

	if err := <-errc; err != nil {
		return err
	}

	err := copyAndFlush(w, cmderr, conn)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func copyAndFlush(w io.Writer, r io.Reader, conn http.Flusher) error {
	br := bufio.NewReader(r)
	for {
		n, err := br.ReadBytes('\n')
		if len(n) > 0 {
			_, err := w.Write(n)
			if err != nil {
				return err
			}
			conn.Flush()
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}
