package usecases

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
)

type FetchGuildMembers struct {
	AlbionService *services.AlbionService
}

func (uc *FetchGuildMembers) Handle() ([]entities.GuildMember, error) {
	return uc.AlbionService.FetchGuildMembers()
}
