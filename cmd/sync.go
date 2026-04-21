package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cwang0126/cursor-synchronizer/internal/fsutil"
	"github.com/cwang0126/cursor-synchronizer/internal/prompts"
)

// syncOptions controls how copyEntries handles overwrite conflicts.
type syncOptions struct {
	// assumeYes forces overwrite without prompting.
	assumeYes bool
}

// copyEntries copies the listed entries from srcCursorRoot to dstCursorRoot,
// prompting on per-file conflicts unless opts.assumeYes is set. It returns
// the list of file paths (relative to .cursor/) that were either written or
// confirmed identical, suitable for inclusion in the manifest.
func copyEntries(srcCursorRoot, dstCursorRoot string, entries []fsutil.Entry, opts syncOptions) ([]string, error) {
	var written []string

	mode := promptMode
	if opts.assumeYes {
		mode = forceAll
	}

	for _, e := range entries {
		files, err := fsutil.CollectFiles(srcCursorRoot, e)
		if err != nil {
			return nil, err
		}
		for _, rel := range files {
			src := filepath.Join(srcCursorRoot, rel)
			dst := filepath.Join(dstCursorRoot, rel)

			if fsutil.Exists(dst) {
				equal, err := fsutil.FilesEqual(src, dst)
				if err != nil {
					return nil, err
				}
				if equal {
					written = append(written, rel)
					continue
				}
				switch mode {
				case skipAll:
					fmt.Fprintf(os.Stderr, "  skip   %s\n", rel)
					written = append(written, rel)
					continue
				case forceAll:
					// fall through to copy
				case promptMode:
					decision, err := prompts.ConfirmOverwrite(rel)
					if err != nil {
						return nil, err
					}
					switch decision {
					case prompts.OverwriteNo:
						fmt.Fprintf(os.Stderr, "  skip   %s\n", rel)
						written = append(written, rel)
						continue
					case prompts.OverwriteSkipAll:
						mode = skipAll
						fmt.Fprintf(os.Stderr, "  skip   %s\n", rel)
						written = append(written, rel)
						continue
					case prompts.OverwriteAll:
						mode = forceAll
					case prompts.OverwriteYes:
						// fall through
					}
				}
			}

			if err := fsutil.CopyFile(src, dst); err != nil {
				return nil, fmt.Errorf("copy %s: %w", rel, err)
			}
			fmt.Fprintf(os.Stderr, "  write  %s\n", rel)
			written = append(written, rel)
		}
	}

	return written, nil
}

type overwriteMode int

const (
	promptMode overwriteMode = iota
	forceAll
	skipAll
)
