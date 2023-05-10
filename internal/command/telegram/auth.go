package telegram

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	indexlog "go.octolab.org/toolset/indexit/internal/log"
	tgsvc "go.octolab.org/toolset/indexit/internal/telegram"
	tgproxy "go.octolab.org/toolset/indexit/internal/telegram/proxy"
)

func authCommand(opt *options) *cobra.Command {
	command := cobra.Command{
		Use:   "auth",
		Short: "Manage Telegram authentication",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(
		authLoginCommand(opt),
		authStatusCommand(opt),
		authLogoutCommand(opt),
	)
	return &command
}

func authLoginCommand(opt *options) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Run interactive Telegram login",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := pathsFromFlags(opt)
			if err != nil {
				return err
			}
			ctx, cancel, err := contextFromFlags(cmd, opt)
			if err != nil {
				return err
			}
			defer cancel()
			client, err := newClient(cmd, paths)
			if err != nil {
				return err
			}
			return tgsvc.Login(ctx, client, newPromptAuth(cmd.InOrStdin(), cmd.ErrOrStderr()))
		},
	}
}

func authStatusCommand(opt *options) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print Telegram session status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := pathsFromFlags(opt)
			if err != nil {
				return err
			}
			if descriptor, derr := tgproxy.FromEnv(); derr != nil {
				return usageErr(derr)
			} else if descriptor != nil {
				indexlog.FromContext(cmd.Context()).Logger.Info("proxy",
					"type", string(descriptor.Type),
					"host", descriptor.Host,
					"port", descriptor.Port)
			}
			if _, err := os.Stat(paths.Session); os.IsNotExist(err) {
				return tgsvc.PrintStatus(cmd.OutOrStdout(), tgsvc.Status{SessionPath: paths.Session})
			} else if err != nil {
				return err
			}
			ctx, cancel, err := contextFromFlags(cmd, opt)
			if err != nil {
				return err
			}
			defer cancel()
			client, err := newClient(cmd, paths)
			if err != nil {
				return err
			}
			status, err := tgsvc.GetStatus(ctx, client, paths.Session)
			if err != nil {
				return err
			}
			return tgsvc.PrintStatus(cmd.OutOrStdout(), status)
		},
	}
}

func authLogoutCommand(opt *options) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Invalidate Telegram session and remove local session file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := pathsFromFlags(opt)
			if err != nil {
				return err
			}
			if _, err := os.Stat(paths.Session); os.IsNotExist(err) {
				return nil
			} else if err != nil {
				return err
			}
			ctx, cancel, err := contextFromFlags(cmd, opt)
			if err != nil {
				return err
			}
			defer cancel()
			client, err := newClient(cmd, paths)
			if err != nil {
				return err
			}
			return tgsvc.Logout(ctx, client, paths.Session)
		},
	}
}

type promptAuth struct {
	in     io.Reader
	out    io.Writer
	reader *bufio.Reader
}

func newPromptAuth(in io.Reader, out io.Writer) auth.UserAuthenticator {
	return &promptAuth{
		in:     in,
		out:    out,
		reader: bufio.NewReader(in),
	}
}

func (p *promptAuth) Phone(ctx context.Context) (string, error) {
	return p.prompt("phone: ")
}

func (p *promptAuth) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	return p.prompt("code: ")
}

func (p *promptAuth) Password(ctx context.Context) (string, error) {
	if file, ok := p.in.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		if _, err := fmt.Fprint(p.out, "2FA password: "); err != nil {
			return "", err
		}
		password, err := term.ReadPassword(int(file.Fd()))
		if _, printErr := fmt.Fprintln(p.out); printErr != nil && err == nil {
			err = printErr
		}
		return string(password), err
	}
	return p.prompt("2FA password: ")
}

func (p *promptAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return fmt.Errorf("sign-up and terms-of-service acceptance are out of scope for this PoC")
}

func (p *promptAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("sign-up is out of scope for this PoC")
}

func (p *promptAuth) prompt(label string) (string, error) {
	if _, err := fmt.Fprint(p.out, label); err != nil {
		return "", err
	}
	line, err := p.reader.ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
