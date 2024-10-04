package main

import "testing"
import "os/exec"

func TestGitHooksSetup(t *testing.T) {
	// git config --local include.path ../.gitconfig
	cmd := exec.Command("git", "config", "--local", "include.path", "../.gitconfig")
	if cmd.Err == nil {
		t.Log("Updated local git config to use .gitconfig")
		return
	}
	output, err := cmd.Output()
	if err == nil {
		t.Errorf("Failed to update local git config: %s\nOutput:\n%s", cmd.Err, string(output))
	} else {
		t.Errorf("Failed to update local git config: %s\nFailed to read git command output: %s", cmd.Err, err)
	}
}
