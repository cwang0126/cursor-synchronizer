package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwang0126/cursor-synchronizer/internal/config"
	"github.com/cwang0126/cursor-synchronizer/internal/fsutil"
	"github.com/cwang0126/cursor-synchronizer/internal/git"
	"github.com/cwang0126/cursor-synchronizer/internal/prompts"
	"github.com/spf13/cobra"
)

var (
	cloneAll    bool
	cloneBranch string
	cloneFolder string
)

var cloneCmd = &cobra.Command{
	Use:   "clone <repo-url> [directory]",
	Short: "Clone .cursor config from a remote repo into a project directory",
	Long: `Shallow-clones the given repo, lets you select which rules/skills/commands
to import, copies them into <directory>/.cursor/, and writes
<directory>/.cursor-sync/{config.yaml,manifest.yaml}.

By default the source layout is auto-detected on the remote, in order:
  - .cursor/{rules,skills,commands}/   (preferred)
  - cursor/{rules,skills,commands}/
  - {rules,skills,commands}/           (at the repo root)

Pass --folder <path> to skip auto-detection and read from a specific folder
relative to the repo root (e.g. --folder configs/cursor). The chosen folder
is recorded in .cursor-sync/config.yaml and reused by later pulls.

If [directory] is omitted, the current folder (.) is used.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runClone,
}

func init() {
	cloneCmd.Flags().BoolVarP(&cloneAll, "all", "a", false, "Skip the multi-select and import everything")
	cloneCmd.Flags().StringVarP(&cloneBranch, "branch", "b", "", "Branch of the remote repo to use (default: try main, then master)")
	cloneCmd.Flags().StringVarP(&cloneFolder, "folder", "f", "", "Remote folder (relative to repo root) containing rules/skills/commands (default: auto-detect .cursor, cursor, or repo root)")
}

func runClone(cmd *cobra.Command, args []string) error {
	repoURL := args[0]
	if !looksLikeRepoURL(repoURL) {
		return fmt.Errorf("first argument must be a repo URL (got %q)", repoURL)
	}
	targetDir := "."
	if len(args) == 2 {
		targetDir = args[1]
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absTarget, 0o755); err != nil {
		return fmt.Errorf("create target dir: %w", err)
	}

	tmp, branch, err := cloneWithBranchFallback(repoURL, cloneBranch)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	srcCursor, err := resolveOrDetectCursorRoot(tmp, cloneFolder)
	if err != nil {
		return err
	}

	entries, err := fsutil.ListCursorEntries(srcCursor)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("remote has no rules/skills/commands to sync")
	}

	var selected []fsutil.Entry
	if cloneAll {
		selected = entries
	} else {
		selected, err = prompts.SelectEntries(entries)
		if err != nil {
			return err
		}
		if len(selected) == 0 {
			fmt.Fprintln(os.Stderr, "Nothing selected; aborting.")
			return nil
		}
	}

	dstCursor := filepath.Join(absTarget, config.CursorDir)
	written, err := copyEntries(srcCursor, dstCursor, selected, syncOptions{})
	if err != nil {
		return err
	}

	if err := config.Save(absTarget, &config.Config{
		Remote: repoURL,
		Branch: branch,
		Folder: cloneFolder,
	}); err != nil {
		return err
	}
	if err := config.SaveManifest(absTarget, &config.Manifest{Entries: written}); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "\nDone. Wrote %d file(s) to %s\n", len(written), dstCursor)
	return nil
}

// cloneWithBranchFallback shallow-clones repoURL. If branch is non-empty it
// is used as-is; otherwise we try "main" first and fall back to "master".
// Returns the temp dir path and the branch that actually succeeded.
func cloneWithBranchFallback(repoURL, branch string) (string, string, error) {
	if branch != "" {
		fmt.Fprintf(os.Stderr, "Cloning %s (branch %s)...\n", repoURL, branch)
		tmp, err := git.ShallowClone(repoURL, branch)
		if err != nil {
			return "", "", err
		}
		return tmp, branch, nil
	}

	candidates := []string{"main", "master"}
	for i, candidate := range candidates {
		fmt.Fprintf(os.Stderr, "Cloning %s (branch %s)...\n", repoURL, candidate)
		tmp, err := git.ShallowClone(repoURL, candidate)
		if err == nil {
			return tmp, candidate, nil
		}
		if i < len(candidates)-1 {
			fmt.Fprintf(os.Stderr, "branch %s not found, trying %s...\n", candidate, candidates[i+1])
		}
	}
	return "", "", fmt.Errorf("could not find branch main or master on remote; pass --branch <name> to use a different one")
}

// looksLikeRepoURL returns true for the URL forms git itself accepts so we
// can detect when the user accidentally passes a directory as the first arg.
func looksLikeRepoURL(s string) bool {
	switch {
	case strings.HasPrefix(s, "http://"),
		strings.HasPrefix(s, "https://"),
		strings.HasPrefix(s, "ssh://"),
		strings.HasPrefix(s, "git://"),
		strings.HasPrefix(s, "file://"),
		strings.HasPrefix(s, "git@"),
		strings.HasSuffix(s, ".git"):
		return true
	}
	return false
}
