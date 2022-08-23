package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCmdCompletion creates a completion command.
func NewCmdCompletion() *cobra.Command {
	validArgs := []string{"bash", "zsh", "fish", "powershell"}
	cmd := &cobra.Command{
		Use:       "completion SHELL",
		Short:     "Generate shell completion",
		Long:      "Generate shell completion. Support: " + strings.Join(validArgs, ", "),
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: validArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := args[0]
			switch sh {
			case "bash":
				err := rootCmd.GenBashCompletion(cmd.OutOrStdout())
				if err != nil {
					return errors.Wrapf(err, "failed to bash completion")
				}
			case "zsh":
				err := rootCmd.GenZshCompletion(cmd.OutOrStdout())
				if err != nil {
					return errors.Wrapf(err, "failed to zsh completion")
				}
			case "fish":
				err := rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
				if err != nil {
					return errors.Wrapf(err, "failed to fish completion")
				}
			case "powershell":
				err := rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
				if err != nil {
					return errors.Wrapf(err, "failed to PowerShell completion")
				}
			default:
				err := showHelp(cmd, args)
				if err != nil {
					return errors.Wrapf(err, "failed to show help")
				}
			}
			return nil
		},
	}

	return cmd
}
