package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Checkout checks out a repo.
func Checkout(b, l string) error {
	cmd := exec.Command("git", "checkout", b)
	cmd.Dir = l

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("checkout %s: Could not run git command", l)
	}

	return nil
}

// Clean cleans a repo.
func Clean(l string) error {
	cmd := exec.Command("git", "clean", "-f")
	cmd.Dir = l

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("clean %s: Could not run git command", l)
	}

	return nil
}

// Clone clones a repo.
func Clone(u, b, l string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", "-b", b, u)
	cmd.Dir = l

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("clone %s: Could not run git command", l)
	}

	return nil
}

// Diff checks a repo for differences.
func Diff(b, l string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-status", "origin/"+b)
	cmd.Dir = l
	bb := new(bytes.Buffer)
	cmd.Stdout = bb

	err := cmd.Run()
	if err != nil {
		return []string{}, fmt.Errorf("diff %s: Could not run git command", l)
	}

	d := bb.String()
	if len(d) < 1 {
		return []string{}, nil
	}

	// Make output pretty.
	d = strings.Replace(d, "A\t", "Adding ", -1)
	d = strings.Replace(d, "C\t", "Copying ", -1)
	d = strings.Replace(d, "D\t", "Deleting ", -1)
	d = strings.Replace(d, "M\t", "Editing ", -1)
	d = strings.Replace(d, "R\t", "Renaming ", -1)
	dl := strings.Split(d, "\n")
	sort.Strings(dl)

	return dl[1:], nil
}

// Fetch fetches a repo.
func Fetch(l string) error {
	cmd := exec.Command("git", "fetch", "--depth", "1")
	cmd.Dir = l

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("fetch %s: Could not run git command", l)
	}

	return nil
}

// Reset resets a repo.
func Reset(b, l string) error {
	cmd := exec.Command("git", "reset", "--hard", "origin/"+b)
	cmd.Dir = l

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("reset %s: Could not run git command", l)
	}

	return nil
}
