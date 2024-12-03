package listeners

import (
	"github.com/bwmarrin/discordgo"
)

// MessageListener escucha mensajes y procesa comandos basados en prefijos.
type MessageListener struct{}

// NewMessageListener crea un nuevo MessageListener.
func NewMessageListener() *MessageListener {
	return &MessageListener{}
}

// HandleMessage procesa los mensajes entrantes.
func (ml *MessageListener) HandleMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Determinar qué comando fue ejecutado
	switch i.ApplicationCommandData().Name {
	case "killboard":
		guild := i.ApplicationCommandData().Options[0].StringValue()
		response := "Registrando killboard para la guild: " + guild

		// Responder a la interacción
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	case "set":
		channelType := i.ApplicationCommandData().Options[0].StringValue()
		response := "Configurando el canal como tipo: " + channelType

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Comando no reconocido.",
			},
		})
	}
}
