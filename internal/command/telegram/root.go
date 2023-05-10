package telegram

import (
	"github.com/spf13/cobra"
)

type options struct {
	sessionPath string
	peersPath   string
	timeout     string
}

func New() *cobra.Command {
	var opt options
	command := cobra.Command{
		Use:   "telegram",
		Short: "Fetch Telegram dialogs and message history",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.PersistentFlags().StringVar(&opt.sessionPath, "session", "", "path to Telegram session file")
	command.PersistentFlags().StringVar(&opt.peersPath, "peer-cache", "", "path to Telegram peer cache file")
	command.PersistentFlags().StringVar(&opt.timeout, "timeout", "0", "wall-clock timeout for Telegram operations")
	command.AddCommand(
		authCommand(&opt),
		fetchCommand(&opt),
		debugCommand(),
	)
	return &command
}
