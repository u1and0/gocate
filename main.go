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
	"io/ioutil"
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
	db  arrayField
	err error
	up  bool
)

type arrayField []string

// arrayField.String sets multiple -f flag
func (i *arrayField) String() string {
	// change this, this is just can example to satisfy the interface
	return "my string representation"
}

// arrayField.Set sets multiple -f flag
func (i *arrayField) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

// db に:が含まれていたら、分割して[]stringに格納
func (sa *arrayField) splitCollon() (sb arrayField) {
	for _, s := range sa {
		if strings.Contains(s, ":") {
			sa := strings.Split(s, ":")
			dbpath = append(sb, da...)
		} else {
			dbpath = append(sb, d)
		}
	}
	return
}

func readOpt() []string {
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Var(&db, "d", "", "Path of locate database file (ex: /path/something.db:/path/another.db)")
	flag.StringVar(&db, "database", "", "Path of locate database file (ex: /path/something.db:/path/another.db)")
	flag.BoolVar(&up, "init", false, "updatedb mode")
	flag.Usage = func() {
		usageTxt := `parallel find files by name

Usage of gocate
	gocate [OPTION]... PATTERN...

-v, -version
	Show version
-d, -database string
	Path of locate database file (ex: /path/something.db:/path/another.db)
-init
	updatedb mode
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
	if len(db) < 1 { // -d option が設定されなかったら
		db = arrayField{os.Getenv("LOCATE_PATH")}
		if len(db) < 1 { // LOCATE_PATHをdbとする
			db = arrayField{DEFAULTDB} // LOCATE_PATH も設定されなかったら DEFAULTDBとする
		} else { // LOCATE_PATHが設定されていたら
			// 2重検索を止めるためにLOCATE_PATHを空にする
			if err := os.Setenv("LOCATE_PATH", ""); err != nil {
				panic(err)
			}
			// 終了時にLOCATE_PATHを戻して終了
			defer func() {
				if err := os.Setenv("LOCATE_PATH", db[0]); err != nil {
					panic(err)
				}
			}()
		}
	}

	db = db.splitCollon()

	// Run updatedb
	if up {
		for _, d := range db {
			dirs, err := ioutil.ReadDir(d)
			if err != nil {
				panic(err)
			}
			for _, d := range dirs {
				com.Wg.Add(1)
				go com.Updatedb(db, d)
			}
			com.Wg.Wait()
		}
		return
	}

	// Run locate
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go cmd.Receiver(c)
	for _, d := range strings.Split(db, ":") {
		com.Wg.Add(1) // カウンタの追加はExec()の外でないとすぐ終わる
		go com.Exec(d, c)
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
