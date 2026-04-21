package git

import (
	"fmt"
	"os"
	"os/exec"
)

// ShallowClone runs `git clone --depth 1 --branch <branch> <url>` into a
// fresh OS temp directory and returns its absolute path.
//
// The caller is responsible for removing the directory when done (typically
// via `defer os.RemoveAll(dir)`).
func ShallowClone(url, branch string) (string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", fmt.Errorf("`git` not found on PATH: install Git to use cursor-sync")
	}

	tmp, err := os.MkdirTemp("", "cursor-sync-*")
	if err != nil {
		return "", err
	}

	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, url, tmp)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(tmp)
		return "", fmt.Errorf("git clone failed: %w", err)
	}
	return tmp, nil
}
