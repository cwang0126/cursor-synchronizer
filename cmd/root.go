package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cwang0126/cursor-synchronizer/internal/banner"
	"github.com/spf13/cobra"
)

const version = "0.3.0"

// rootUsageTemplate is cobra's default usage template with two tweaks:
//   - the Available Commands rows render "name|alias1|alias2" via cmdLabel
//     instead of just "name"
//   - group/help-topic branches are stripped since we don't use them
//
// The indentation around {{...}} tags is significant — keep template output
// pixel-identical to cobra's default layout.
const rootUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad (cmdLabel .) .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

const tagline = "cursor-sync keeps your local .cursor config in sync with a remote Git repository."

// bareHeader is printed only for the bare `cursor-sync` invocation, directly
// below the banner. It's intentionally kept off the --help output.
var bareHeader = banner.Cyan + tagline + banner.Reset + "\n\033[2mVersion " + version + "\033[0m"

var rootCmd = &cobra.Command{
	Use:   "cursor-sync",
	Short: "Sync .cursor rules/skills/commands from a remote git repo",
	Example: `  cursor-sync clone <repo-url> [directory]
  cursor-sync clone <repo-url> --branch [branch-name] --folder [folder-path]`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	// Accept any positional args at the root so bare invocations like
	// `cursor-sync <url>` reach runRoot instead of cobra's default
	// "unknown command" error.
	Args: cobra.ArbitraryArgs,
	Run:  runRoot,
}

// Execute runs the root cobra command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.EnableCommandSorting = false
	cobra.AddTemplateFunc("cmdLabel", cmdLabel)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetUsageTemplate(rootUsageTemplate)

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)
}

// cmdLabel returns the string shown in the "Available Commands" column of
// help output: the command's primary name, plus any aliases joined with "|".
func cmdLabel(c *cobra.Command) string {
	if len(c.Aliases) == 0 {
		return c.Name()
	}
	return c.Name() + "|" + strings.Join(c.Aliases, "|")
}

// runRoot handles bare invocations of `cursor-sync`.
//
//   - No args: print the banner plus the full --help output, so users
//     discover every available command without having to know the --help flag.
//   - A repo-URL-looking arg: suggest the likely-intended `cursor-sync clone`.
//   - Any other unknown arg: short error pointing at --help.
func runRoot(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		banner.Print(os.Stdout)
		fmt.Fprintln(os.Stdout, bareHeader)
		fmt.Fprintln(os.Stdout)
		_ = cmd.Usage()
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown command or argument: %q\n", args[0])
	if looksLikeRepoURL(args[0]) {
		fmt.Fprintf(os.Stderr, "\nDid you mean:\n  cursor-sync clone %s\n\n", args[0])
	}
	fmt.Fprintln(os.Stderr, "Run `cursor-sync --help` to see available commands.")
}
