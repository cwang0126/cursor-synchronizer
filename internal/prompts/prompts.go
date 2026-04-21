package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cwang0126/cursor-synchronizer/internal/fsutil"
)

// selectAllLabel is the meta-option shown first in the multi-select.
const selectAllLabel = "[Select All]"

// SelectEntries shows a multi-select prompt for the given entries with a
// "Select All" item as the first option. Arrow keys navigate, space toggles,
// Enter confirms.
//
// If the user picks "[Select All]", every entry is returned regardless of
// other selections.
func SelectEntries(entries []fsutil.Entry) ([]fsutil.Entry, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	options := make([]string, 0, len(entries)+1)
	options = append(options, selectAllLabel)
	byLabel := make(map[string]fsutil.Entry, len(entries))
	for _, e := range entries {
		label := e.RelPath()
		if e.IsDir {
			label += "/"
		}
		options = append(options, label)
		byLabel[label] = e
	}

	var picked []string
	q := &survey.MultiSelect{
		Message:  "Select rules/skills/commands to sync:",
		Options:  options,
		PageSize: 15,
	}
	if err := survey.AskOne(q, &picked); err != nil {
		return nil, err
	}

	for _, p := range picked {
		if p == selectAllLabel {
			return entries, nil
		}
	}

	out := make([]fsutil.Entry, 0, len(picked))
	for _, p := range picked {
		if e, ok := byLabel[p]; ok {
			out = append(out, e)
		}
	}
	return out, nil
}

// OverwriteDecision is the result of an overwrite prompt.
type OverwriteDecision int

const (
	OverwriteNo OverwriteDecision = iota
	OverwriteYes
	OverwriteAll
	OverwriteSkipAll
)

// ConfirmOverwrite asks the user whether to overwrite a file.
//
// Accepts: y (yes), N (no, default), a (yes-to-all), s (skip-all).
func ConfirmOverwrite(path string) (OverwriteDecision, error) {
	var ans string
	q := &survey.Input{
		Message: fmt.Sprintf("Overwrite %s? [y/N/a/s]", path),
		Default: "N",
		Help:    "y=yes  N=no (default)  a=yes to all  s=skip all",
	}
	if err := survey.AskOne(q, &ans); err != nil {
		return OverwriteNo, err
	}
	switch ans {
	case "y", "Y":
		return OverwriteYes, nil
	case "a", "A":
		return OverwriteAll, nil
	case "s", "S":
		return OverwriteSkipAll, nil
	default:
		return OverwriteNo, nil
	}
}
