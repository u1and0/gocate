package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

const (
	// BENCH : for benchmark test
	BENCH = false
)

// Command : Command executer
type Command struct {
	// Args : Search keyword and option
	Args []string
	// output : Storing directory for updatedb database
	Output string
}

// Receiver : channel receiver
func Receiver(ch <-chan string) {
	for {
		s, ok := <-ch
		if !ok {
			break
		}
		if BENCH {
			continue // Ignore print for benchmark test
		}
		fmt.Println(s)
	}
}

// Locate : locate command generator
func (c *Command) Locate(dir string) *exec.Cmd {
	// locate command option read after -- from command line
	opt := append([]string{"-d", dir}, remove(c.Args, "--")...)
	command := exec.Command("locate", opt...)
	return command
}

// Run : locate command executer write to channel
func Run(c exec.Cmd, ch chan string) error {
	stdout, err := c.StdoutPipe()
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	// Command execute
	if err := c.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		if s := scanner.Text(); s != "" {
			// time.Sleep(1 * time.Millisecond)  // [test]順序守らないことのマーカー
			ch <- s
		}
	}

	// Command Error handling
	// go func() {
	b, _ := io.ReadAll(stderr)
	if err = errors.New(string(b)); err != nil {
		return err
	}
	// }()

	if err := c.Wait(); err != nil {
		return err
	}
	return err
}

// remove specified string from string array
func remove(ss []string, s string) []string {
	for i, v := range ss {
		if v == s {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}
