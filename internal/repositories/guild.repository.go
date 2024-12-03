package repositories

import "albion-killbot/internal/entities"

type GuildRepository interface {
	FetchGuildMembers() ([]entities.GuildMember, error)
}
