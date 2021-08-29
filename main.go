// Usage:
//		gocate [-d path] [--database=path] [--version] [--help] pattern...
// For benchmark test, const BENCH turns true then run below
//		$ go test -bench
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	// BENCH : Benchmark test flag
	BENCH bool = false
	// VERSION : Show version flag
	VERSION string = "v0.1.0"
)

var (
	showVersion bool
	// for normalLocate test default value
	db   string = "./test/var.db:./test/etc.db:./test/usr.db"
	word string = ".*pacman.*proto"
)

func receiver(ch <-chan string) {
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

func main() {
	// Read option
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()
	if showVersion {
		fmt.Println("gocate version:", VERSION)
		return // Exit with version info
	}

	db = os.Getenv("LOCATE_PATH")
	word = strings.Join(flag.Args(), ".*")

	// Run goroutine
	var wg sync.WaitGroup // カウンタを宣言
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go receiver(c)
	for _, o := range strings.Split(db, ":") {
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
					// time.Sleep(1 * time.Millisecond)  // [test]順序守らないことのマーカー
					c <- s
				}
			}
		}(o)
	}
	wg.Wait() // カウンタが0になるまでブロック
}

// Nomral locate command for benchmark
func normalLocate() {
	b, _ := exec.Command("locate", "-i", "-d", db, "--regex", word).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		if BENCH {
			continue
		}
		fmt.Println(o)
	}
}

/*
$ go test -bench .
goos: linuxbench .
goarch: amd64
pkg: speedtest/src/github.com/u1and0/gocate
cpu: Intel(R) Core(TM) i5-8400 CPU @ 2.80GHz
BenchmarkNormalLocate-6               13          89004938 ns/op
BenchmarkParallelLocate-6             12          88001858 ns/op
PASS
ok      speedtest/src/github.com/u1and0/gocate  2.401s


普通のlocateより1msecくらい勝つようになった
normalLocateのほうに--regexオプションがないからだった
*/
