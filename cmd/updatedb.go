package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func (c *Command) Updatedb(root string, d fs.FileInfo) {
	defer c.Wg.Done()
	if d.IsDir() {
		cur, _ := os.Getwd()
		opt := []string{
			"-U",
			fmt.Sprintf("%s/%s", root, d.Name()),
			"--output",
			fmt.Sprintf("%s/%s/%s.db", cur, ".gocate", d.Name()),
			// `pwd`/.gocate is temp directory
		}
		command := exec.Command("updatedb", opt...)
		fmt.Println(command)
		command.Run()
	}
}
