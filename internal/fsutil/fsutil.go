package fsutil

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CursorSubdirs are the top-level groupings inside .cursor/ that we sync.
var CursorSubdirs = []string{"rules", "skills", "commands"}

// CursorRootCandidates are the directories, in priority order, that
// DetectCursorRoot probes when looking for the source layout in a remote
// repo. The empty string means "the repo root itself".
var CursorRootCandidates = []string{".cursor", "cursor", ""}

// DetectCursorRoot finds which directory inside repoRoot holds the
// rules/skills/commands groupings. It tries, in order:
//
//  1. repoRoot/.cursor/
//  2. repoRoot/cursor/
//  3. repoRoot/           (groupings at the repo root)
//
// A candidate wins if it contains at least one of CursorSubdirs. Returns
// "", false when no candidate matches.
func DetectCursorRoot(repoRoot string) (string, bool) {
	for _, candidate := range CursorRootCandidates {
		root := filepath.Join(repoRoot, candidate)
		if !Exists(root) {
			continue
		}
		for _, sub := range CursorSubdirs {
			if Exists(filepath.Join(root, sub)) {
				return root, true
			}
		}
	}
	return "", false
}

// Entry is one top-level item under .cursor/<group>/, e.g. a single rule
// file or a skill folder.
type Entry struct {
	// Group is one of "rules", "skills", "commands".
	Group string
	// Name is the top-level entry name within the group (file or dir).
	Name string
	// IsDir is true when the entry is a directory (e.g. a skill folder).
	IsDir bool
}

// RelPath returns the entry path relative to .cursor/, e.g. "rules/foo.mdc".
func (e Entry) RelPath() string {
	return filepath.Join(e.Group, e.Name)
}

// ListCursorEntries enumerates top-level entries under cursorRoot/{rules,skills,commands}.
//
// Missing subdirectories are silently skipped (a remote may only define a
// subset of groupings).
func ListCursorEntries(cursorRoot string) ([]Entry, error) {
	var out []Entry
	for _, group := range CursorSubdirs {
		dir := filepath.Join(cursorRoot, group)
		items, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", dir, err)
		}
		for _, it := range items {
			out = append(out, Entry{
				Group: group,
				Name:  it.Name(),
				IsDir: it.IsDir(),
			})
		}
	}
	return out, nil
}

// CollectFiles walks an entry rooted at cursorRoot and returns the list of
// file paths relative to cursorRoot.
func CollectFiles(cursorRoot string, e Entry) ([]string, error) {
	root := filepath.Join(cursorRoot, e.RelPath())
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{e.RelPath()}, nil
	}
	var out []string
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(cursorRoot, path)
		if err != nil {
			return err
		}
		out = append(out, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CopyFile copies src to dst, creating parent directories as needed.
// It overwrites dst if it already exists.
func CopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp := dst + ".tmp-cursor-sync"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}

// Exists reports whether path exists (file or directory).
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FilesEqual reports whether two files have identical contents using SHA-256.
// Returns false (no error) if either file is missing.
func FilesEqual(a, b string) (bool, error) {
	ha, err := hashFile(a)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	hb, err := hashFile(b)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	for i := range ha {
		if ha[i] != hb[i] {
			return false, nil
		}
	}
	return true, nil
}

func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
