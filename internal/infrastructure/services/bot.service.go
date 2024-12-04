package services

import (
	"albion-killbot/internal/entities"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type BotService struct {
	Session *discordgo.Session
}

// NewBotService crea una nueva instancia del servicio del bot
func NewBotService(session *discordgo.Session) *BotService {
	return &BotService{Session: session}
}

// RegisterCommand registra un comando en Discord
func (b *BotService) RegisterCommand(command entities.Command) error {
	// Crear el comando para la API de Discord
	discordCommand := &discordgo.ApplicationCommand{
		Name:        command.Name,
		Description: command.Description,
		Options:     convertOptions(command.Options),
	}

	// Registrar el comando con la API de Discord
	_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", discordCommand)
	if err != nil {
		return fmt.Errorf("error al registrar el comando %s: %v", command.Name, err)
	}

	log.Printf("Comando registrado: %s", command.Name)
	return nil
}

// OpenSession abre la sesión del bot
func (b *BotService) OpenSession() error {
	if b.Session == nil {
		return fmt.Errorf("sesión de Discord no inicializada")
	}
	// Verificar que la sesión esté abierta
	if b.Session.Open() != nil {
		return fmt.Errorf("no se pudo abrir la sesión de Discord")
	}
	return nil
}

// CloseSession cierra la sesión del bot
func (b *BotService) CloseSession() error {
	if b.Session == nil {
		return fmt.Errorf("sesión de Discord no inicializada")
	}
	return b.Session.Close()
}

// ConvertOptions convierte opciones personalizadas a las de DiscordGo
func convertOptions(options []entities.CommandOption) []*discordgo.ApplicationCommandOption {
	var botOptions []*discordgo.ApplicationCommandOption
	for _, opt := range options {
		botOpt := &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionType(opt.Type),
			Name:        opt.Name,
			Description: opt.Description,
			Required:    opt.Required,
		}

		if len(opt.Choices) > 0 {
			for _, choice := range opt.Choices {
				botOpt.Choices = append(botOpt.Choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  choice.Name,
					Value: choice.Value,
				})
			}
		}

		botOptions = append(botOptions, botOpt)
	}
	return botOptions
}

func (b *BotService) SendInteractionToServer(channelId string, message string) error {
	_, err := b.Session.ChannelMessageSend(channelId, message)
	if err != nil {
		log.Printf("Error al enviar interacción al servidor: %v", err)
		return err
	}
	log.Printf("Interacción enviada al servidor: %s, mensaje: %s", channelId, message)
	return nil
}
