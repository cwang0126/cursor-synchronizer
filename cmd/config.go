package cmd

import (
	"fmt"
	"os"

	"github.com/cwang0126/cursor-synchronizer/internal/config"
	"github.com/spf13/cobra"
)

var (
	configShow string
	configSet  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or update fields in .cursor-sync/config.yaml",
	Long: `Examples:
  cursor-sync config --show remote
  cursor-sync config --set remote https://github.com/owner/repo.git
  cursor-sync config --set branch main
  cursor-sync config --set folder configs/cursor`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfig,
}

func init() {
	configCmd.Flags().StringVar(&configShow, "show", "", "Show the value of a config field (e.g. remote, branch, folder)")
	configCmd.Flags().StringVar(&configSet, "set", "", "Set a config field; pass the new value as a positional argument")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if configShow == "" && configSet == "" {
		return fmt.Errorf("must pass --show <field> or --set <field> <value>")
	}
	if configShow != "" && configSet != "" {
		return fmt.Errorf("--show and --set are mutually exclusive")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if configShow != "" {
		cfg, err := config.Load(cwd)
		if err != nil {
			return err
		}
		val, err := getField(cfg, configShow)
		if err != nil {
			return err
		}
		fmt.Println(val)
		return nil
	}

	if len(args) != 1 {
		return fmt.Errorf("--set %s requires a value argument", configSet)
	}
	value := args[0]

	cfg, err := config.Load(cwd)
	if err != nil {
		return err
	}
	if err := setField(cfg, configSet, value); err != nil {
		return err
	}
	if err := config.Save(cwd, cfg); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Updated %s = %s\n", configSet, value)
	return nil
}

func getField(c *config.Config, field string) (string, error) {
	switch field {
	case "remote":
		return c.Remote, nil
	case "branch":
		return c.Branch, nil
	case "folder":
		return c.Folder, nil
	default:
		return "", fmt.Errorf("unknown field %q (supported: remote, branch, folder)", field)
	}
}

func setField(c *config.Config, field, value string) error {
	switch field {
	case "remote":
		c.Remote = value
	case "branch":
		c.Branch = value
	case "folder":
		c.Folder = value
	default:
		return fmt.Errorf("unknown field %q (supported: remote, branch, folder)", field)
	}
	return nil
}
