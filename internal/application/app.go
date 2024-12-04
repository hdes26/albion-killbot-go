package application

import (
	"albion-killbot/internal/infrastructure/db"
	dbrepositories "albion-killbot/internal/infrastructure/db/repositories"
	"albion-killbot/internal/listeners"
	"context"
	"fmt"
	"log"
)

type App struct {
	bot          *Bot
	dbClient     *db.MongoDBClient
	ChannelRepo  *dbrepositories.ChannelRepository
	KillListener *listeners.KillListener
}

func NewApp(botToken string, dbClient *db.MongoDBClient) *App {

	channelCollection := dbClient.GetDatabase("albion").Collection("channels")
	channelRepo := &dbrepositories.ChannelRepository{DB: channelCollection}

	bot := NewBot(botToken, channelRepo)

	killListener := listeners.NewKillListener(bot.Session, channelRepo)

	return &App{
		bot:          bot,
		dbClient:     dbClient,
		ChannelRepo:  channelRepo,
		KillListener: killListener,
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("App started")

	// Iniciar el listener que maneja los eventos
	go a.KillListener.Start(ctx, a.ChannelRepo)

	// Iniciar el bot en un goroutine para que se ejecute de manera concurrente
	errCh := make(chan error)

	// Ejecutar el bot en segundo plano
	go func() {
		if err := a.bot.Run(ctx); err != nil {
			errCh <- fmt.Errorf("error al ejecutar el bot: %v", err)
		}
	}()

	// Esperar que el contexto termine (normalmente, esto sería una cancelación o un cierre de la app)
	<-ctx.Done()

	// Esperar por errores del bot
	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
