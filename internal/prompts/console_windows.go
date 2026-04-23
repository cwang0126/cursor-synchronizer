//go:build windows

package prompts

import (
	"os"

	"golang.org/x/sys/windows"
)

// disableVTInput clears ENABLE_VIRTUAL_TERMINAL_INPUT on stdin for the
// duration of a survey prompt and returns a closure that restores the
// previous console mode.
//
// survey/v2's Windows rune reader (terminal/runereader_windows.go) only
// clears ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT and
// reads keys via ReadConsoleInputW, mapping VK_UP / VK_DOWN / VK_LEFT /
// VK_RIGHT to its arrow-key constants. When ENABLE_VIRTUAL_TERMINAL_INPUT
// is set (e.g. Windows Terminal after `git` has written progress output to
// the same console), arrow presses arrive as raw VT100 escape sequences
// (ESC `[` `A` / `B` / `C` / `D`) instead of arrow KEY_EVENT records, and
// survey echoes the trailing `[A` / `[B` bytes into the filter buffer.
func disableVTInput() func() {
	h := windows.Handle(os.Stdin.Fd())
	var prev uint32
	if err := windows.GetConsoleMode(h, &prev); err != nil {
		return func() {}
	}
	if prev&windows.ENABLE_VIRTUAL_TERMINAL_INPUT == 0 {
		return func() {}
	}
	if err := windows.SetConsoleMode(h, prev&^windows.ENABLE_VIRTUAL_TERMINAL_INPUT); err != nil {
		return func() {}
	}
	return func() {
		_ = windows.SetConsoleMode(h, prev)
	}
}
