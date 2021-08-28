// Run
// $ go test -v
package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Opt : locate command options
type Opt struct {
	Dir  string
	Word string
}

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	wg.Add(2)             // カウンタの初期化
	opts := []Opt{
		{"./test/bin.db", "ls"},
		{"./test/usr.db", "ing.hpp"},
	}
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
		// time.Sleep(1 * time.Microsecond)
		fmt.Println(o)
	}
	wg.Done() // カウンタを減算
}

/* Nomral locate command for benchmark
func locateBin0() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/bin.db", "ls").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
}

func locateUsr0() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "ing.hpp").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
}
*/
