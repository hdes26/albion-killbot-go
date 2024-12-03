package usecases

import (
	"albion-killbot/internal/entities"

	"github.com/bwmarrin/discordgo"
)

// GenerateKillEventEmbeds genera los embeds para un evento de combate
func GenerateKillEventEmbeds(event entities.Event) []*discordgo.MessageEmbed {
	killerEmbed := &discordgo.MessageEmbed{
		Color: 0x0099FF,
		Title: "[" + event.KillerGuild + "] " + event.KillerName + " killed [" + event.VictimGuild + "] " + event.VictimName,
		Image: &discordgo.MessageEmbedImage{URL: "attachment://kill" + event.EventID + ".jpg"},
	}

	victimEmbed := &discordgo.MessageEmbed{
		Color:  0xFF0000,
		Image:  &discordgo.MessageEmbedImage{URL: "attachment://inventory" + event.EventID + ".jpg"},
		Footer: &discordgo.MessageEmbedFooter{Text: "Powered by ache !"},
	}

	return []*discordgo.MessageEmbed{killerEmbed, victimEmbed}
}
