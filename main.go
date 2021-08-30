// Build:
// $ go build
// Usage:
//		gocate [-d path] [--database=path] [--version] [--help] PATTERN... -- [LOCATE OPTION]
//		$ ./gocate -d $(find test -name '*.db' | paste -sd:) -- -i --regex fstab
// For benchmark test, const BENCH turns true then run below
//		$ go test -bench
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	cmd "gocate/cmd"
)

const (
	// BENCH : Benchmark test flag
	BENCH bool = false
	// VERSION : Show version flag
	VERSION string = "v0.1.1r"
	// DEFAULTDB : Default locate search path
	DEFAULTDB string = "/var/lib/mlocate/mlocate.db"
)

var (
	com cmd.Command
	// locate command path
	showVersion bool
	// for normalLocate test default value
	db  string
	err error
	// word for test
)

func readOpt() []string {
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
		os.Exit(0) // Exit with version info
	}
	return flag.Args() // options + search word
}

func main() {
	// Check locate command
	com.Exe, err = exec.LookPath("locate")
	if err != nil {
		panic(err)
	}
	com.Args = readOpt()
	com.Wg = sync.WaitGroup{} // カウンタを宣言

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
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go cmd.Receiver(c)
	for _, d := range strings.Split(db, ":") {
		com.Wg.Add(1) // カウンタの追加
		com.Dir = d
		go com.Exec(c)
	}
	com.Wg.Wait() // カウンタが0になるまでブロック
}

// Nomral locate command for benchmark
func normalLocate(args []string) {
	b, _ := exec.Command("locate", args...).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		if BENCH {
			continue
		}
		fmt.Println(o)
	}
}
