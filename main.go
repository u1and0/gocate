// Run
// $ go test -v
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Opt : locate command options
type Opt struct {
	Dir  string
	Word string
}

const (
	opts  string = "./test/var.db:./test/etc.db:./test/usr.db"
	word  string = "pacman"
	bench bool   = true
)

func receiver(ch <-chan string) {
	for {
		s, ok := <-ch
		if !ok {
			break
		}
		if bench {
			continue
		}
		fmt.Println(s)
	}
}

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go receiver(c)
	for _, o := range strings.Split(opts, ":") {
		wg.Add(1) // カウンタの追加
		go func(o string) {
			defer wg.Done() // go func抜けるときにカウンタを減算
			cmd := exec.Command("locate", "-i", "-d", o, "--regex", word)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cmd.Start()

			scanner := bufio.NewScanner(stdout)

			for scanner.Scan() {
				if s := scanner.Text(); s != "" {
					c <- s
				}
			}
		}(o)
	}
	wg.Wait() // カウンタが0になるまでブロック
}

// // Locate excutes locate command
// func Locate(o Opt, wg *sync.WaitGroup) {
// 	b, _ := exec.Command("locate", "-i", "-d", o.Dir, o.Word).Output()
// 	out := strings.Split(string(b), "\n")
// 	for _, o := range out {
// 		time.Sleep(1 * time.Microsecond)
// 		fmt.Println(o)
// 	}
// 	wg.Done() // カウンタを減算
// }

// Nomral locate command for benchmark
func normalLocate() {
	b, _ := exec.Command("locate", "-i", "-d", opts, word).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		if bench {
			continue
		}
		fmt.Println(o)
	}
}

/*
$ go test -bench .
goos: linux
goarch: amd64
pkg: speedtest/src/github.com/u1and0/gocate
cpu: Intel(R) Core(TM) i5-8400 CPU @ 2.80GHz
BenchmarkNormalLocate-6               56          21427959 ns/op
BenchmarkParallelLocate-6             37          33087714 ns/op
PASS
ok      speedtest/src/github.com/u1and0/gocate  2.482s

普通のlocateの方が1.5倍速いまでに近づいた
*/
