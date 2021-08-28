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

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	wg.Add(2)             // カウンタの初期化
	go locateBin(&wg)
	go locateUsr(&wg)
	wg.Wait() // カウンタが0になるまでブロック
}

func locateBin(wg *sync.WaitGroup) {
	b, _ := exec.Command("locate", "-i", "-d", "./test/bin.db", "ls").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
	wg.Done() // カウンタを減算
}

func locateUsr(wg *sync.WaitGroup) {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "ing.hpp").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
	wg.Done() // カウンタを減算
}

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
