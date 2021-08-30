// Build:
// $ go build
// Usage:
//		gocate [-d path] [--database=path] [--version] [--help] PATTERN... -- [LOCATE OPTION]
//		$ ./gocate -d $(find test -name '*.db' | paste -sd:) -- -i --regex fstab
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
	// DEFAULTDB : Default locate search path
	DEFAULTDB string = "/var/lib/mlocate/mlocate.db"
)

var (
	showVersion bool
	// for normalLocate test default value
	db string
	// word for test
	word = []string{
		"-i",
		"-d",
		"./test/var.db:./test/etc.db:./test/usr.db",
		"--regex",
		".*pacman.*proto",
	}
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
	flag.StringVar(&db, "d", "", "Path of locate database file (ex: /path/something.db:/path/another.db)")
	flag.StringVar(&db, "database", "", "Path of locate database file (ex: /path/something.db:/path/another.db)")
	flag.Usage = func() {
		usageTxt := `parallel find files by name

Usage of gocate
	gocate [OPTION]... PATTERN...

-v, -version
	Show version
-d, -database string
	Path of locate database file (ex: /path/something.db:/path/another.db)
-- [OPTION]...
	locate command option`
		fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
	}
	flag.Parse()
	if showVersion {
		fmt.Println("gocate version:", VERSION)
		return // Exit with version info
	}
	word = flag.Args() // options + search word

	// db 優先順位
	// -d PATH > LOCATE_PATH > /var/lib/mlocate/mlocate.db
	if db == "" { // -d option が設定されなかったら
		if db = os.Getenv("LOCATE_PATH"); db == "" { // LOCATE_PATHをdbとする
			db = DEFAULTDB // LOCATE_PATH も設定されなかったら DEFAULTDBとする
		} else { // LOCATE_PATHが設定されていたら
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
		}
	}

	// Run goroutine
	var wg sync.WaitGroup // カウンタを宣言
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go receiver(c)
	for _, o := range strings.Split(db, ":") {
		wg.Add(1) // カウンタの追加
		go func(o string) {
			defer wg.Done() // go func抜けるときにカウンタを減算

			// locate command option read after -- from command line
			opt := append([]string{"-d", o}, word...)
			cmd := exec.Command("locate", opt...)
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
	b, _ := exec.Command("locate", word...).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		if BENCH {
			continue
		}
		fmt.Println(o)
	}
}
