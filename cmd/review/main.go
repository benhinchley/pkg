// Command review simplifies pushing changes to gerrit
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	branch string
	remote string
)

func main() {
	flag.StringVar(&branch, "branch", "master", "branch to push refs for")
	flag.StringVar(&remote, "remote", "origin", "gerrit remote to push refs to")
	flag.Parse()

	refs := fmt.Sprintf("HEAD:refs/for/%s", branch)
	cmd := exec.Command("git", "push", remote, refs)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to run \"git push %s HEAD:/refs/for/%s\"\n", remote, branch)
		os.Exit(1)
	}
}
