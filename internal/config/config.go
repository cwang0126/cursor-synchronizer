package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

const (
	// Dir is the per-project config directory.
	Dir = ".cursor-sync"
	// CursorDir is the per-project Cursor config directory we sync into.
	CursorDir = ".cursor"

	configFile   = "config.yaml"
	manifestFile = "manifest.yaml"
)

// Config is the contents of .cursor-sync/config.yaml.
type Config struct {
	Remote string `yaml:"remote"`
	Branch string `yaml:"branch"`
}

// Manifest is the contents of .cursor-sync/manifest.yaml.
//
// Entries are paths relative to the .cursor/ directory, e.g.
// "rules/karpathy-guidelines.mdc" or "skills/deslop".
type Manifest struct {
	Entries []string `yaml:"entries"`
}

// ConfigPath returns the absolute path to the config file under projectDir.
func ConfigPath(projectDir string) string {
	return filepath.Join(projectDir, Dir, configFile)
}

// ManifestPath returns the absolute path to the manifest file under projectDir.
func ManifestPath(projectDir string) string {
	return filepath.Join(projectDir, Dir, manifestFile)
}

// Load reads .cursor-sync/config.yaml from projectDir.
//
// Returns a wrapped error if the file is missing so callers can produce a
// friendly "run `cursor-sync clone` first" message.
func Load(projectDir string) (*Config, error) {
	path := ConfigPath(projectDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("no %s found in %s: run `cursor-sync clone` first", filepath.Join(Dir, configFile), projectDir)
		}
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Branch == "" {
		c.Branch = "master"
	}
	return &c, nil
}

// Save writes .cursor-sync/config.yaml under projectDir, creating the
// parent directory if needed.
func Save(projectDir string, c *Config) error {
	if err := os.MkdirAll(filepath.Join(projectDir, Dir), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(projectDir), data, 0o644)
}

// LoadManifest reads .cursor-sync/manifest.yaml. Missing file returns an
// empty manifest (not an error) so first-run callers don't need a special
// path.
func LoadManifest(projectDir string) (*Manifest, error) {
	path := ManifestPath(projectDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Manifest{}, nil
		}
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &m, nil
}

// SaveManifest writes the manifest, deduping and sorting entries for stable diffs.
func SaveManifest(projectDir string, m *Manifest) error {
	if err := os.MkdirAll(filepath.Join(projectDir, Dir), 0o755); err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(m.Entries))
	out := make([]string, 0, len(m.Entries))
	for _, e := range m.Entries {
		if _, ok := seen[e]; ok {
			continue
		}
		seen[e] = struct{}{}
		out = append(out, e)
	}
	sort.Strings(out)
	data, err := yaml.Marshal(&Manifest{Entries: out})
	if err != nil {
		return err
	}
	return os.WriteFile(ManifestPath(projectDir), data, 0o644)
}
