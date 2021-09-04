package cmd

import (
	"os/exec"
	"testing"
)

func TestUpdatedb(t *testing.T) {
	com := Command{Gocatedbpath: "./.gocate"}
	expected := exec.Command("updatedb", "-U", "/usr/bin", "--output", ".gocate/usr_bin.db")
	actual := com.Updatedb("/usr/bin")

	if expected != actual {
		t.Fatalf("\ngot:  %v\nwant: %v", actual, expected)
	}
}
