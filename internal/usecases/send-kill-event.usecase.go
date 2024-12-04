package usecases

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// SendInteractionUseCase es el caso de uso para enviar interacciones al servidor de Discord
type SendKillEventUseCase struct {
	BotService *services.BotService
}

// NewSendInteractionUseCase crea una nueva instancia del caso de uso
func NewSendKillEeventUseCase(botService *services.BotService) *SendKillEventUseCase {
	return &SendKillEventUseCase{
		BotService: botService,
	}
}

// Handle envía la interacción al servidor de Discord
func (uc *SendKillEventUseCase) Handle(channelId string, event entities.PlayerKill) error {
	embed := GenerateKillEventEmbeds(event)
	err := uc.BotService.SendInteractionToServer(channelId, embed)
	if err != nil {
		return fmt.Errorf("error al enviar interacción: %v", err)
	}
	return nil
}

func GenerateKillEventEmbeds(kill entities.PlayerKill) []*discordgo.MessageEmbed {
	killerEmbed := &discordgo.MessageEmbed{
		Color: 0x0099FF,
		Title: "[" + kill.Killer.GuildName + "] " + kill.Killer.Name + " killed [" + kill.Victim.GuildName + "] " + kill.Victim.Name,
		/* Image: &discordgo.MessageEmbedImage{URL: "attachment://kill" + event.EventID + ".jpg"}, */
	}

	victimEmbed := &discordgo.MessageEmbed{
		Color: 0xFF0000,
		/* Image:  &discordgo.MessageEmbedImage{URL: "attachment://inventory" + kill.EventId + ".jpg"}, */
		Footer: &discordgo.MessageEmbedFooter{Text: "Powered by ache !"},
	}

	return []*discordgo.MessageEmbed{killerEmbed, victimEmbed}
}
