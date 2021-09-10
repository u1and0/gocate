package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	// BENCH : for benchmark test
	BENCH = false
)

// Command : Command executer
type Command struct {
	// Args : Search keyword and option
	Args []string
	// Wg : Waiting group for goroutine
	Wg   sync.WaitGroup
	// Gocatedbpath : Storing directory for updatedb database
	Gocatedbpath string
}

// Receiver : channel receiver
func Receiver(ch <-chan string) {
	for {
		s, ok := <-ch
		if !ok {
			break
		}
		if BENCH {
			continue
		}
		fmt.Println(s)
	}
}

// Locate : locate command executer
func (c *Command) Locate(dir string, ch chan string) {
	defer c.Wg.Done() // go func抜けるときにカウンタを減算

	// locate command option read after -- from command line
	opt := append([]string{"-d", dir}, c.Args...)
	command := exec.Command("locate", opt...)
	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command.Start()

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		if s := scanner.Text(); s != "" {
			// time.Sleep(1 * time.Millisecond)  // [test]順序守らないことのマーカー
			ch <- s
		}
	}
}
