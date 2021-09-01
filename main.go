// Build:
// $ go build
// Usage:
//		gocate [-d path] [--database=path] [--version] [--help] PATTERN... -- [LOCATE OPTION]
//		$ ./gocate -d $(find test -name '*.db' | paste -sd:) -- -i --regex 'lib.*id$'
// For benchmark test, const BENCH turns true then run below
//		$ go test -bench
package main

import (
	"flag"
	"fmt"
	"io/fs"
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
	db string
	// updatedb mode flag
	up bool
	// updatedb path
	updb = arrayField{"/"}
	//
)

type arrayField []string

type UsageText struct {
	showVersion string
	db          string
	up          string
	updb        string
}

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

// // db に:が含まれていたら、分割して[]stringに格納
// func (sa *arrayField) splitCollon() (sb arrayField) {
// 	for _, s := range sa {
// 		if strings.Contains(s, ":") {
// 			sa := strings.Split(s, ":")
// 			dbpath = append(sb, da...)
// 		} else {
// 			dbpath = append(sb, d)
// 		}
// 	}
// 	return
// }

func readOpt() []string {
	usage := UsageText{
		showVersion: "Show version",
		db:          "Path of locate database file (ex: /path/something.db:/path/another.db)",
		up:          "updatedb mode",
		updb:        "Store only results of scanning the file system subtree rooted at PATH  to  the  generated  database.",
	}
	flag.BoolVar(&showVersion, "v", false, usage.showVersion)
	flag.BoolVar(&showVersion, "version", false, usage.showVersion)
	flag.StringVar(&db, "d", "", usage.db)
	flag.StringVar(&db, "database", "", usage.db)
	flag.BoolVar(&up, "init", false, usage.up)
	flag.Var(&updb, "U", usage.updb)
	flag.Var(&updb, "database-root", usage.updb)
	flag.Usage = func() {
		usageTxt := fmt.Sprintf(`parallel find files by name

Usage of gocate
	gocate [OPTION]... PATTERN...

-v, -version
	%s
-d, -database string
	%s
-init
	%s
-U, -database-root
	%s
-- [OPTION]...
	locate command option`, usage.showVersion, usage.db, usage.up, usage.updb)
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
	for _, c := range []string{"locate", "updatedb"} {
		_, err := exec.LookPath(c)
		if err != nil {
			panic(err)
		}
	}
	com.Gocatedbpath = "./.gocate" // "/var/lib/mlocate"
	com.Args = readOpt()
	com.Wg = sync.WaitGroup{} // カウンタを宣言

	// db 優先順位
	// -d PATH > LOCATE_PATH > /var/lib/mlocate/mlocate.db
	if len(db) < 1 { // -d option が設定されなかったら
		db = os.Getenv("LOCATE_PATH")
		if len(db) < 1 { // LOCATE_PATHをdbとする
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

	// db = db.splitCollon()

	// Run updatedb
	if up {
		// <= -U /usr -U /etc
		for _, pairent := range updb { // updb = arrayField{"/usr", "/etc"}
			dirs, err := ioutil.ReadDir(pairent) // =>fs.FileInfo{ /usr/lib, /usr/bin, ... /etc/pacman.d}
			if err != nil {
				panic(err)
			}
			for _, d := range dirs {
				com.Wg.Add(1)
				go func(d fs.FileInfo) {
					if err := com.Updatedb(pairent, d); err != nil {
						panic(err)
					}
				}(d)
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
		/* arrayField db はパスを複数持っている
		 * `gocate -d /usr -d /etc:/var` として走らせた場合
		 * "/usr", "/etc:/var" コロンで区切られた場合は、
		 * そのままlocateに渡して1データベースとして検索する
		 */
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
