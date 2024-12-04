package application

import (
	"albion-killbot/internal/infrastructure/db"
	dbrepositories "albion-killbot/internal/infrastructure/db/repositories"
	"albion-killbot/internal/infrastructure/services"
	"albion-killbot/internal/listeners"
	"albion-killbot/internal/usecases"
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Bot es el controlador que orquesta la ejecución del bot
type Bot struct {
	Session                    *discordgo.Session
	BotRegisterCommandsUseCase *usecases.BotRegisterCommandsUseCase // Referencia al caso de uso
	SendKillEventUseCase       *usecases.SendKillEventUseCase
	ChannelRepo                *dbrepositories.ChannelRepository
}

func NewBot(botToken string, dbClient *db.MongoDBClient) *Bot {
	// Crear la sesión de Discord
	sess, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil
	}
	botService := &services.BotService{
		Session: sess,
	}

	if err := botService.OpenSession(); err != nil {
		log.Println("Error abriendo la sesión de Discord:", err)
		return nil
	}

	// Obtener la colección de canales de la base de datos
	channelCollection := dbClient.GetDatabase("albion").Collection("channels")

	// Crear el ChannelRepository
	channelRepo := dbrepositories.NewChannelRepository(channelCollection)

	// Crear el caso de uso para registrar comandos
	botRegisterCommandsUseCase := &usecases.BotRegisterCommandsUseCase{
		BotService: botService, // Aquí pasas botService directamente, no hace falta hacer &botService
	}

	// Crear y devolver el bot
	return &Bot{
		Session:                    sess,
		BotRegisterCommandsUseCase: botRegisterCommandsUseCase,
		ChannelRepo:                channelRepo, // Pasamos ChannelRepo al bot
	}
}

// Run ejecuta el bot, registrando comandos y manteniendo la sesión activa
func (b *Bot) Run(ctx context.Context) error {
	// Registrar comandos del bot
	err := b.BotRegisterCommandsUseCase.Handle() // No pasa la sesión ahora
	if err != nil {
		log.Fatalf("Error al registrar comandos: %v", err)
		return err
	}

	log.Println("Bot is online")

	/* Bot listeners */
	messageListener := listeners.NewMessageListener(b.ChannelRepo) // Le pasamos ChannelRepo

	// Registrar el listener en la sesión
	b.Session.AddHandler(messageListener.HandleMessage)

	b.Session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	return nil
}

func (b *Bot) SendInteractionToServer(channelId string, message string) error {
	return b.SendKillEventUseCase.Handle(channelId, message)
}
