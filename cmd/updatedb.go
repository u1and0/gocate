package cmd

import (
	"fmt"
	"io/fs"
)

func (c *Command) Updatedb(d fs.FileInfo) {
	defer c.Wg.Done()
	if d.IsDir() {
		fmt.Println(d.Name())
	}
}
