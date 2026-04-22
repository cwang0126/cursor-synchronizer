package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwang0126/cursor-synchronizer/internal/config"
	"github.com/cwang0126/cursor-synchronizer/internal/fsutil"
	"github.com/cwang0126/cursor-synchronizer/internal/git"
	"github.com/spf13/cobra"
)

var pullYes bool

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull the latest .cursor entries previously synced into this project",
	Long: `Re-syncs only the entries already tracked in .cursor-sync/manifest.yaml.
On per-file conflicts the user is prompted (y/N/a/s); pass --yes to overwrite all.`,
	Args: cobra.NoArgs,
	RunE: runPull,
}

func init() {
	pullCmd.Flags().BoolVarP(&pullYes, "yes", "y", false, "Overwrite all conflicting files without prompting")
}

func runPull(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return err
	}
	manifest, err := config.LoadManifest(cwd)
	if err != nil {
		return err
	}
	if len(manifest.Entries) == 0 {
		return fmt.Errorf("manifest is empty: nothing to pull. Run `cursor-sync clone` first.")
	}

	tracked := topLevelEntriesFromManifest(manifest.Entries)
	if len(tracked) == 0 {
		return fmt.Errorf("manifest contains no recognizable entries")
	}

	fmt.Fprintf(os.Stderr, "Pulling %s (branch %s)...\n", cfg.Remote, cfg.Branch)
	tmp, err := git.ShallowClone(cfg.Remote, cfg.Branch)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	srcCursor, ok := fsutil.DetectCursorRoot(tmp)
	if !ok {
		return fmt.Errorf("remote repo has no .cursor/, cursor/, or root-level rules/skills/commands directories")
	}

	remoteEntries, err := fsutil.ListCursorEntries(srcCursor)
	if err != nil {
		return err
	}

	var toSync []fsutil.Entry
	missing := make([]string, 0)
	for key := range tracked {
		var found *fsutil.Entry
		for i := range remoteEntries {
			if remoteEntries[i].RelPath() == key {
				found = &remoteEntries[i]
				break
			}
		}
		if found == nil {
			missing = append(missing, key)
			continue
		}
		toSync = append(toSync, *found)
	}

	for _, m := range missing {
		fmt.Fprintf(os.Stderr, "  warn   %s no longer exists on remote (kept locally)\n", m)
	}

	if len(toSync) == 0 {
		fmt.Fprintln(os.Stderr, "Nothing to pull.")
		return nil
	}

	dstCursor := filepath.Join(cwd, config.CursorDir)
	written, err := copyEntries(srcCursor, dstCursor, toSync, syncOptions{assumeYes: pullYes})
	if err != nil {
		return err
	}

	if err := config.SaveManifest(cwd, &config.Manifest{Entries: written}); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "\nDone. Synced %d file(s).\n", len(written))
	return nil
}

// topLevelEntriesFromManifest derives the set of group/name keys (e.g.
// "skills/deslop") from the flat list of file paths in the manifest.
func topLevelEntriesFromManifest(entries []string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, e := range entries {
		parts := strings.SplitN(filepath.ToSlash(e), "/", 3)
		if len(parts) < 2 {
			continue
		}
		out[parts[0]+"/"+parts[1]] = struct{}{}
	}
	return out
}
