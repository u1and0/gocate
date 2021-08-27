// Run
// $ go test -v
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// nothing to do
}

func locateBin() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/bin.db", "ls").Output()
	fmt.Printf("%v", string(b))
}

func locateUsr() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "ing.hpp").Output()
	fmt.Printf("%v", string(b))
}

func normalLocate() {
	locateUsr()
	locateBin()
}

func parallelLocate() {
	go locateUsr()
	go locateBin()
}
