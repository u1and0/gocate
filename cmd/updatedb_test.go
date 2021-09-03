package cmd

import (
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestUpdatedb(t *testing.T) {
	com := Command{Gocatedbpath: "./gocate"}
	expected := exec.Command("updatedb", "-U", "/usr/bin", "--output", "./gocate/usr_bin.db")
	dirs, _ := ioutil.ReadDir("/usr")
	f := dirs[0]
	actual := com.Updatedb(f)

	if expected != actual {
		t.Fatalf("\ngot:  %v\nwant: %v", actual, expected)
	}
}
