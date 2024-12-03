package listeners

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/usecases"
	"fmt"

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
	fmt.Println(i.ChannelID)
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Determinar qué comando fue ejecutado
	switch i.ApplicationCommandData().Name {
	case "killboard":
		guild := i.ApplicationCommandData().Options[0].StringValue()

		event := entities.Event{
			EventID:     "12345",
			KillerGuild: guild,
			KillerName:  "KillerPlayer",
			VictimGuild: "VictimGuild",
			VictimName:  "VictimPlayer",
		}

		embeds := usecases.GenerateKillEventEmbeds(event)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embeds,
			},
		})
		if err != nil {
			// Manejar error en la respuesta
			println("Error enviando la respuesta:", err.Error())
		}

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
