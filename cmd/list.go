package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cwang0126/cursor-synchronizer/internal/config"
	"github.com/cwang0126/cursor-synchronizer/internal/fsutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List local .cursor entries and tag them as [remote] or [local]",
	Long: `Lists rules/skills/commands under ./.cursor/ in the current directory.
Entries whose top-level name is recorded in .cursor-sync/manifest.yaml are
tagged [remote]; user-added entries are tagged [local]. Offline; no network call.`,
	Args: cobra.NoArgs,
	RunE: runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cursorRoot := filepath.Join(cwd, config.CursorDir)
	if !fsutil.Exists(cursorRoot) {
		return fmt.Errorf("no %s/ directory in %s", config.CursorDir, cwd)
	}

	manifest, err := config.LoadManifest(cwd)
	if err != nil {
		return err
	}
	tracked := topLevelEntriesFromManifest(manifest.Entries)

	entries, err := fsutil.ListCursorEntries(cursorRoot)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("(no entries under .cursor/)")
		return nil
	}

	currentGroup := ""
	for _, e := range entries {
		if e.Group != currentGroup {
			fmt.Printf("\n%s/\n", e.Group)
			currentGroup = e.Group
		}
		tag := "[local]"
		if _, ok := tracked[e.RelPath()]; ok {
			tag = "[remote]"
		}
		suffix := ""
		if e.IsDir {
			suffix = "/"
		}
		fmt.Printf("  %-9s %s%s\n", tag, e.Name, suffix)
	}
	return nil
}
