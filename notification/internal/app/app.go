package app

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/voronovsg/rocket-factory/notification/internal/config"
	"github.com/voronovsg/rocket-factory/platform/pkg/closer"
	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 3)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runOrderPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("order paid consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runOrderAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("order assembled consumer crashed: %v", err)
		}
	}()

	go func() {
		a.runTelegramBot(ctx)
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initTelegramBot,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initTelegramBot(_ context.Context) error {
	telegramBot := a.diContainer.TelegramBot()

	telegramBot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			logger.Info(ctx, "Telegram bot activation", zap.Int64("chat_id", update.Message.Chat.ID))

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "🚀 Notification Bot активирован! Теперь вы будете получать уведомления о статусе ваших заказов.",
			})
			if err != nil {
				logger.Error(ctx, "Failed to send activation message", zap.Error(err))
			}
		})

	return nil
}

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 NotificationService Order Paid consumer running")

	err := a.diContainer.OrderPaidConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runOrderAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 NotificationService Order Assembled consumer running")

	err := a.diContainer.OrderAssembledConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runTelegramBot(ctx context.Context) {
	logger.Info(ctx, "🤖 NotificationService Telegram bot running")

	telegramBot := a.diContainer.TelegramBot()
	telegramBot.Start(ctx)
}
