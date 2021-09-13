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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	cmd "github.com/u1and0/gocate/cmd"
)

const (
	// BENCH : Benchmark test flag
	BENCH = false
	// VERSION : Show version flag
	VERSION = "v0.2.2r"
	// DEFAULTDB : Default locate search path
	DEFAULTDB = "/var/lib/mlocate/"
)

var (
	// com command structure
	com cmd.Command
	// locate command path
	showVersion bool
	// for normalLocate test default value
	db string
	// updatedb mode flag
	up bool
	// updatedb path
	updb = arrayField{}
	// output directory for updatedb
	output string
	// if true, do not run the updatedb script
	dryrun bool
)

type arrayField []string

type usageText struct {
	showVersion string
	db          string
	up          string
	updb        string
	output      string
	dryrun      string
}

// arrayField.String sets multiple -f flag
func (a *arrayField) String() string {
	// change this, this is just can example to satisfy the interface
	return "my string representation"
}

// arrayField.Set sets multiple -f flag
func (a *arrayField) Set(value string) error {
	*a = append(*a, strings.TrimSpace(value))
	return nil
}

// Dbpath : directory names for indexing using -U option
func (a *arrayField) Dbpath() (dd []string) {
	for _, pairent := range *a { // a = arrayField{"/usr", "/etc"}
		dirs, err := ioutil.ReadDir(pairent) // => fs.FileInfo{ lib, bin, ... }
		if err != nil {
			panic(err)
		}
		ft := cmd.FileTree{Pairent: pairent, Dirs: dirs} // fss = /usr/bin /usr/lib ... ( []fs.FileInfo )
		dd = ft.DirectoryFilter(dd)
	}
	return
}

func flagParse() []string {
	usage := usageText{
		showVersion: "Show version",
		db:          "Path of locate database file (default: /var/lib/mlocate)",
		up:          "updatedb mode",
		updb:        "Store only results of scanning the file system subtree rooted at PATH  to  the  generated  database.",
		output:      "Write the database to DIRECTORY instead of using the default database directory.",
		dryrun:      "Just print command, do NOT run updatedb command.",
	}
	flag.BoolVar(&showVersion, "v", false, usage.showVersion)
	flag.BoolVar(&showVersion, "version", false, usage.showVersion)
	flag.StringVar(&db, "d", "/var/lib/mlocate", usage.db)
	flag.StringVar(&db, "database", "/var/lib/mlocate", usage.db)
	flag.BoolVar(&up, "init", false, usage.up)
	flag.Var(&updb, "U", usage.updb)
	flag.Var(&updb, "database-root", usage.updb)
	flag.StringVar(&output, "o", DEFAULTDB, usage.output)
	flag.StringVar(&output, "output", DEFAULTDB, usage.output)
	flag.BoolVar(&dryrun, "dryrun", false, usage.dryrun)
	flag.Usage = func() {
		usageTxt := fmt.Sprintf(`parallel find files by name

Usage of gocate
	gocate [OPTION]... PATTERN...

-v, -version
	%s
-d, -database DIRECTORY
	%s
-init
	%s
-U, -database-root DIRECTORY
	%s
-o, -output DIRECTORY
	%s
-dryrun
	%s
-- [OPTION]...
	locate or updatedb command option`,
			usage.showVersion,
			usage.db,
			usage.up,
			usage.updb,
			usage.output,
			usage.dryrun,
		)
		fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
	}
	flag.Parse()
	if showVersion {
		fmt.Println("gocate version:", VERSION)
		os.Exit(0) // Exit with version info
	}
	if len(updb) < 1 { // updb default value
		updb = arrayField{"/"}
	}
	return flag.Args() // options + search word
}

func main() {
	// Check locate command
	for _, c := range []string{"locate", "updatedb"} {
		if _, err := exec.LookPath(c); err != nil {
			panic(err)
		}
	}
	com.Args = flagParse() // 先に実行しないとoutputとかのフラグ読み込まれない
	com.Output = output
	com.Wg = sync.WaitGroup{} // カウンタを宣言

	// // db 優先順位
	// // -d PATH > LOCATE_PATH > /var/lib/mlocate/mlocate.db
	// if len(db) < 1 { // -d option が設定されなかったら
	// 	db = os.Getenv("LOCATE_PATH")
	// 	if len(db) < 1 { // LOCATE_PATHをdbとする
	// 		db = DEFAULTDB // LOCATE_PATH も設定されなかったら DEFAULTDBとする
	// 	} else { // LOCATE_PATHが設定されていたら
	// 		// 2重検索を止めるためにLOCATE_PATHを空にする
	// 		if err := os.Setenv("LOCATE_PATH", ""); err != nil {
	// 			panic(err)
	// 		}
	// 		// 終了時にLOCATE_PATHを戻して終了
	// 		defer func() {
	// 			if err := os.Setenv("LOCATE_PATH", db); err != nil {
	// 				panic(err)
	// 			}
	// 		}()
	// 	}
	// }

	// Run updatedb
	if up { // <= $ gocate -init -U /usr -U /etc
		for _, dir := range updb.Dbpath() { // => /usr/bin /usr/lib ...
			com.Wg.Add(1)
			go func(d string) {
				defer com.Wg.Done()
				c := com.Updatedb(d)
				fmt.Println(c)
				if !dryrun {
					if err := c.Run(); err != nil {
						panic(err)
					}
				}
			}(dir)
		}
		com.Wg.Wait()
		os.Exit(0)
	}

	// Run locate
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go cmd.Receiver(c)
	dbs, err := filepath.Glob(db + "/*.db")
	if err != nil {
		panic(nil)
	}
	fmt.Println(dbs)
	for _, d := range dbs {
		// for _, d := range strings.Split(db, ":") {
		/* arrayField db はパスを複数持っている
		 * `gocate -d /usr -d /etc:/var` として走らせた場合
		 * "/usr", "/etc:/var" コロンで区切られた場合は、
		 * そのままlocateに渡して1データベースとして検索する
		 */
		com.Wg.Add(1) // カウンタの追加はExec()の外でないとすぐ終わる
		go com.Locate(d, c)
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
