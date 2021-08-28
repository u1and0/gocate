// Run
// $ go test -v
package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Opt : locate command options
type Opt struct {
	Dir  string
	Word string
}

var opts = []Opt{
	{"./test/bin.db", "ls"},
	{"./test/usr.db", "ing.hpp"},
}

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	wg.Add(2)             // カウンタの初期化
	for _, o := range opts {
		go Locate(o, &wg)
	}
	wg.Wait() // カウンタが0になるまでブロック
}

// Locate excutes locate command
func Locate(o Opt, wg *sync.WaitGroup) {
	b, _ := exec.Command("locate", "-i", "-d", o.Dir, o.Word).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Println(o)
	}
	wg.Done() // カウンタを減算
}

// Nomral locate command for benchmark
func locateBin0() {
	b, _ := exec.Command("locate", "-i", "-d", opts[0].Dir, opts[0].Word).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Println(o)
	}
}

func locateUsr0() {
	b, _ := exec.Command("locate", "-i", "-d", opts[1].Dir, opts[1].Word).Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Println(o)
	}
}

/*
$ go test -bench .
goos: linux
goarch: amd64
pkg: speedtest/src/github.com/u1and0/gocate
cpu: Intel(R) Core(TM) i5-8400 CPU @ 2.80GHz
BenchmarkNormalLocate-6               38          32892696 ns/op
BenchmarkParallelLocate-6             42          28137809 ns/op
PASS
ok      speedtest/src/github.com/u1and0/gocate  2.506s
*/
