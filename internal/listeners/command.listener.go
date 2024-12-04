package listeners

import (
	"albion-killbot/internal/entities"
	dbrepositories "albion-killbot/internal/infrastructure/db/repositories"
	"albion-killbot/internal/infrastructure/services"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// MessageListener escucha mensajes y procesa comandos basados en prefijos.
type MessageListener struct {
	AlbionService services.AlbionService
	ChannelRepo   *dbrepositories.ChannelRepository // Repositorio para guardar el canal

}

// NewMessageListener crea un nuevo MessageListener.
func NewMessageListener(channelRepo *dbrepositories.ChannelRepository) *MessageListener {
	return &MessageListener{
		ChannelRepo: channelRepo,
	}
}

// HandleMessage procesa los mensajes entrantes.
func (ml *MessageListener) HandleMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Determinar qué comando fue ejecutado
	switch i.ApplicationCommandData().Name {
	case "killboard":
		exampleEmbed := &discordgo.MessageEmbed{
			Color:       0x0099FF,
			Title:       "Killboard",
			Description: "successfully registered on the channel 💥killboard. Now you need to configure which channels you want to receive notifications on. Use the /killboard set command on the desired channel to receive notifications..",
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg"},
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Enjoy!",
				IconURL: "https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg",
			},
		}

		// Responder con el embed generado
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{exampleEmbed},
			},
		})

		if err != nil {
			// Manejar error en la respuesta
			println("Error enviando la respuesta:", err.Error())
		}

		//TODO:
		channelFound, err := ml.ChannelRepo.FindByChannelID(i.ChannelID) // capturamos ambos valores: canal y error
		if err != nil {
			log.Printf("Error buscando canal por ID: %v", err)
			return
		}

		if channelFound == nil {
			log.Println("Canal no encontrado")
			return
		}

		guildName := i.ApplicationCommandData().Options[0].StringValue()
		guild, err := ml.AlbionService.FetchGuildByName(guildName)

		if err != nil {
			// Maneja el error de no encontrar el guild
			fmt.Println("Guild no encontrado:", guildName)
			return
		}

		// Usar guild directamente, sin desreferenciar
		channel := &entities.Channel{
			ChannelID: i.ChannelID,
			Channel:   guild.Name,
			Type:      nil,
			Active:    nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Guild:     &guild, // Aquí pasas directamente el valor de guild, sin usar *
		}

		_, err = ml.ChannelRepo.Save(channel)
		if err != nil {
			// Maneja el error al guardar el canal
			fmt.Println("Error al guardar el canal:", err)
		}

	case "set":
		// Crear el embed con la respuesta
		exampleEmbed := &discordgo.MessageEmbed{
			Color:       0x0099FF,
			Title:       "Killboard",
			Description: fmt.Sprintf("The channel: %s has been configured to receive Kill notifications. If you haven't already, use the '/killboard guild' command to register your guild.", i.ChannelID),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: "https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg"},
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Enjoy!",
				IconURL: "https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg",
			},
		}

		// Responder con el embed generado
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{exampleEmbed},
			},
		})

		if err != nil {
			// Manejar error en la respuesta
			println("Error enviando la respuesta:", err.Error())
		}
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Comando no reconocido.",
			},
		})
	}
}
