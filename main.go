// Bug
//		Wirte stdout to same result 3 times
// Usage:
//		gocate [-d path] [--database=path] [--version] [--help] pattern...
//		$ LOCATE_PATH=$(find test -name '*.db' | paste -sd:) go run main.go pacman proto
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
	// 2重検索を止めるためにLOCATE_PATHを空にする
	if err := os.Setenv("LOCATE_PATH", ""); err != nil {
		panic(err)
	}
	// 終了時にLOCATE_PATHを戻して終了
	defer func() {
		if err := os.Setenv("LOCATE_PATH", db); err != nil {
			panic(err)
		}
	}()

	// 検索ワード
	word = strings.Join(flag.Args(), " ")

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
