package telegram

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/gotd/td/bin"
	gotd "github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"

	"go.octolab.org/toolset/indexit/internal/buildinfo"
)

type API interface {
	ContactsResolveUsername(context.Context, *tg.ContactsResolveUsernameRequest) (*tg.ContactsResolvedPeer, error)
	MessagesGetDialogs(context.Context, *tg.MessagesGetDialogsRequest) (tg.MessagesDialogsClass, error)
	MessagesGetHistory(context.Context, *tg.MessagesGetHistoryRequest) (tg.MessagesMessagesClass, error)
	MessagesGetReplies(context.Context, *tg.MessagesGetRepliesRequest) (tg.MessagesMessagesClass, error)
	AuthLogOut(context.Context) (*tg.AuthLoggedOut, error)
}

type Runner interface {
	Run(context.Context, func(context.Context, API, *auth.Client) error) error
}

type Client struct {
	client    *gotd.Client
	log       *slog.Logger
	heartbeat time.Duration
}

// ClientOptions carries optional knobs for NewClient. Zero value is valid.
type ClientOptions struct {
	Resolver  dcs.Resolver
	Logger    *slog.Logger // for indexit-level lifecycle markers
	GotdLog   *zap.Logger  // for gotd internals (Nop unless -vvv)
	Heartbeat time.Duration
}

func NewClient(apiID int, apiHash, sessionPath string, opts ClientOptions) *Client {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	gotdLog := opts.GotdLog
	if gotdLog == nil {
		gotdLog = zap.NewNop()
	}

	middlewares := []gotd.Middleware{}
	if opts.Heartbeat > 0 {
		middlewares = append(middlewares, &heartbeatMiddleware{
			log:      logger,
			interval: opts.Heartbeat,
		})
	}

	gotdOpts := gotd.Options{
		NoUpdates:      true,
		Logger:         gotdLog,
		SessionStorage: &gotd.FileSessionStorage{Path: sessionPath},
		Middlewares:    middlewares,
		Device: gotd.DeviceConfig{
			DeviceModel:    "indexit",
			SystemVersion:  runtime.GOOS + "/" + runtime.GOARCH,
			AppVersion:     buildinfo.Display(),
			SystemLangCode: "en",
			LangCode:       "en",
		},
	}
	if opts.Resolver != nil {
		gotdOpts.Resolver = opts.Resolver
	}
	return &Client{
		client:    gotd.NewClient(apiID, apiHash, gotdOpts),
		log:       logger,
		heartbeat: opts.Heartbeat,
	}
}

func (c *Client) Run(ctx context.Context, fn func(context.Context, API, *auth.Client) error) error {
	c.log.Info("telegram: connecting")
	start := time.Now()
	connected := make(chan struct{})

	if c.heartbeat > 0 {
		go func() {
			ticker := time.NewTicker(c.heartbeat)
			defer ticker.Stop()
			for {
				select {
				case <-connected:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					c.log.Info("telegram: still connecting",
						"elapsed", time.Since(start).Round(time.Second))
				}
			}
		}()
	}

	return c.client.Run(ctx, func(ctx context.Context) error {
		close(connected)
		c.log.Info("telegram: connected", "elapsed", time.Since(start).Round(time.Second))
		return fn(ctx, c.client.API(), c.client.Auth())
	})
}

// heartbeatMiddleware logs a "still waiting" line if a single RPC exceeds the
// configured interval. The interval acts as both threshold and tick.
type heartbeatMiddleware struct {
	log      *slog.Logger
	interval time.Duration
}

type tlNamer interface{ TypeName() string }

func methodName(input bin.Encoder) string {
	if n, ok := input.(tlNamer); ok {
		return n.TypeName()
	}
	return "unknown"
}

func (h *heartbeatMiddleware) Handle(next tg.Invoker) gotd.InvokeFunc {
	return func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
		name := methodName(input)
		start := time.Now()
		done := make(chan struct{})

		go func() {
			ticker := time.NewTicker(h.interval)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					h.log.Info("telegram: still waiting",
						"method", name,
						"elapsed", time.Since(start).Round(time.Second))
				}
			}
		}()

		err := next.Invoke(ctx, input, output)
		close(done)
		return err
	}
}
