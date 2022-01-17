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
	VERSION = "v0.3.1r"
	// DEFAULTDB : Default locate search path
	DEFAULTDB = "/var/lib/mlocate/"
)

var (
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
		ft := cmd.FileTree{Pairent: pairent, Dirs: dirs} // ft = /usr/bin /usr/lib ... ( []fs.FileInfo )
		dd = ft.DirectoryFilter(dd)
	}
	return
}

func flagParse() cmd.Command {
	var (
		showVersion bool
		usage       = usageText{
			showVersion: "Show version",
			db:          "Path of locate database directory (default: /var/lib/mlocate)",
			up:          "updatedb mode",
			updb:        "Store only results of scanning the file system subtree rooted at PATH  to  the  generated  database.",
			output:      "Write the database to DIRECTORY instead of using the default database directory. (default: /var/lib/mlocate)",
			dryrun:      "Just print command, do NOT run updatedb command.",
		}
	)
	flag.BoolVar(&showVersion, "v", false, usage.showVersion)
	flag.BoolVar(&showVersion, "version", false, usage.showVersion)
	flag.StringVar(&db, "d", DEFAULTDB, usage.db)
	flag.StringVar(&db, "database", DEFAULTDB, usage.db)
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
	com := cmd.Command{
		Args:   flag.Args(), // options + search word
		Output: output,
	}
	return com
}

func main() {
	var (
		wg  = sync.WaitGroup{} // カウンタを宣言
		com = flagParse()
	)

	// Check locate command
	for _, c := range []string{"locate", "updatedb"} {
		if _, err := exec.LookPath(c); err != nil {
			panic(err)
		}
	}

	// Run updatedb
	if up { // <= $ gocate -init -U /usr -U /etc
		for _, dir := range updb.Dbpath() { // => /usr/bin /usr/lib ...
			wg.Add(1)
			go func(d string) {
				defer wg.Done()
				c := com.Updatedb(d)
				fmt.Println(c)
				if !dryrun {
					if err := c.Run(); err != nil {
						fmt.Printf("%v", err)
						os.Exit(1)
					}
				}
			}(dir)
		}
		wg.Wait()
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
	for _, d := range dbs {
		wg.Add(1) // カウンタの追加はLocate()の外でないとすぐ終わる
		go func(d string, ch chan string) {
			defer wg.Done() // go func抜けるときにカウンタを減算
			c := com.Locate(d)
			if !dryrun {
				if err := cmd.Run(*c, ch); err != nil {
					fmt.Printf("%v", err)
					os.Exit(1)
				}
			} else {
				fmt.Println(c)
			}
		}(d, c)
	}
	wg.Wait() // カウンタが0になるまでブロック
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
