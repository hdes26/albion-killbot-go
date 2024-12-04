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

	switch i.ApplicationCommandData().Name {
	case "killboard":
		handleKillboardCommand(s, i, ml)
	case "set":
		handleSetCommand(s, i, ml)
	default:
		respondWithMessage(s, i, "Comando no reconocido.")
	}
}

func handleKillboardCommand(s *discordgo.Session, i *discordgo.InteractionCreate, ml *MessageListener) {
	embed := createEmbed(
		"Killboard",
		"successfully registered on the channel 💥killboard. Configure which channels to receive notifications using the `/killboard set` command.",
		"https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg",
	)

	respondWithEmbed(s, i, embed)

	channelID := i.ChannelID

	guildName := i.ApplicationCommandData().Options[0].StringValue()

	// Buscar si el canal ya existe en la base de datos
	exist, err := ml.ChannelRepo.FindByChannelID(channelID)
	if err != nil {
		log.Printf("Error buscando canal por ID: %v", err)
		return
	}

	if exist == nil {
		guild, err := ml.AlbionService.FetchGuildByName(guildName)
		if err != nil {
			log.Printf("Guild no encontrado: %v", err)
			return
		}

		newChannel := &entities.Channel{
			ChannelID: channelID,
			Type:      nil,
			Guild:     &guild,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Active:    boolPtr(false),
		}

		if _, err := ml.ChannelRepo.Save(newChannel); err != nil {
			log.Printf("Error al guardar el canal: %v", err)
		}
		return
	}
	if exist.Guild != nil {
		return
	}

	guild, err := ml.AlbionService.FetchGuildByName(guildName)
	if err != nil {
		log.Printf("Error al buscar el guild: %v", err)
		return
	}
	exist.Guild = &guild
	exist.Active = boolPtr(true)

	exist.UpdatedAt = time.Now()
	if _, err := ml.ChannelRepo.Update(exist); err != nil {
		log.Printf("Error al actualizar el canal: %v", err)
	}
}

func handleSetCommand(s *discordgo.Session, i *discordgo.InteractionCreate, ml *MessageListener) {
	embed := createEmbed(
		"Killboard",
		fmt.Sprintf("The channel: %s has been configured to receive Kill notifications. Use `/killboard guild` to register your guild.", i.ChannelID),
		"https://i.ibb.co/6wSQ18j/logo-albion-bot.jpg",
	)

	respondWithEmbed(s, i, embed)

	channelID := i.ChannelID
	channel, err := s.State.Channel(channelID)
	channelType := i.ApplicationCommandData().Options[0].StringValue()

	// Buscar si el canal ya existe en la base de datos
	exist, err := ml.ChannelRepo.FindByChannelID(channelID)
	if err != nil {
		log.Printf("Error buscando canal por ID: %v", err)
		return
	}

	if exist == nil {
		newChannel := &entities.Channel{
			ChannelID: channelID,
			Channel:   &channel.Name,
			Type:      &channelType,
			Guild:     nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Active:    boolPtr(false),
		}

		if _, err := ml.ChannelRepo.Save(newChannel); err != nil {
			log.Printf("Error al guardar el canal: %v", err)
		}
		return
	}

	if exist.Type != nil {
		return
	}

	exist.Type = &channelType
	exist.Channel = &channel.Name
	exist.Active = boolPtr(true)

	exist.UpdatedAt = time.Now()
	if _, err := ml.ChannelRepo.Update(exist); err != nil {
		log.Printf("Error al actualizar el canal: %v", err)
	}
}

func createEmbed(title, description, thumbnail string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0x0099FF,
		Title:       title,
		Description: description,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: thumbnail},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Enjoy!",
			IconURL: thumbnail,
		},
	}
}

func respondWithEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		log.Printf("Error enviando respuesta: %v", err)
	}
}

func respondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Printf("Error enviando respuesta: %v", err)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
