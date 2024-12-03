package application

import (
	"albion-killbot/internal/infrastructure/services"
	"albion-killbot/internal/listeners"
	"albion-killbot/internal/usecases"
	"context"
	"fmt"
	"log"
)

type App struct {
	KillListener *listeners.KillListener
	bot          *Bot
}

func NewApp(botToken string) *App {
	// Inicializamos los servicios y repositorios
	albionService := &services.AlbionService{}

	// Inicializamos los casos de uso
	fetchPlayerKills := usecases.FetchPlayerKills{AlbionService: albionService}
	fetchGuildMembers := usecases.FetchGuildMembers{AlbionService: albionService}

	// Inicializamos el bot con el token
	bot := NewBot(botToken) // Aquí se pasa el token

	// Inicializamos el listener
	killListener := &listeners.KillListener{
		FetchPlayerKills:  fetchPlayerKills,
		FetchGuildMembers: fetchGuildMembers,
	}

	return &App{
		KillListener: killListener,
		bot:          bot,
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println("App started")

	// Iniciar el listener que maneja los eventos
	go a.KillListener.Start(ctx)

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
