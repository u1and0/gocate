package cmd

import (
	"fmt"
	"io/fs"
	"os/exec"
)

func (c *Command) Updatedb(root string, f fs.FileInfo) error {
	defer c.Wg.Done()
	if !f.IsDir() {
		return fmt.Errorf("do not must be included file %s", f.Name())
	}
	opt := []string{
		"-U",
		fmt.Sprintf("%s/%s", root, f.Name()),
		"--output",
		fmt.Sprintf("%s/%s.db", c.Gocatedbpath, f.Name()),
		// `pwd`/.gocate is temp directory
	}
	command := exec.Command("updatedb", opt...)
	fmt.Println(command)
	// command.Run()
	return nil
}
