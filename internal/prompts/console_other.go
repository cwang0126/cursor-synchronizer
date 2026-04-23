//go:build !windows

package prompts

// disableVTInput is a no-op on non-Windows platforms; the Windows-only
// console-mode dance is documented in console_windows.go.
func disableVTInput() func() { return func() {} }
