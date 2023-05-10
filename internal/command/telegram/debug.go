package telegram

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"go.octolab.org/toolset/indexit/internal/telegram/uid"
)

func debugCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "debug",
		Short: "Debug Telegram helpers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(debugUIDCommand())
	return &command
}

func debugUIDCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uid <value>",
		Short: "Parse a Telegram UID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := uid.Parse(args[0])
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(ref)
		},
	}
}
