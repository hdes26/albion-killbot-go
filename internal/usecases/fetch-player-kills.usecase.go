package usecases

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
)

type FetchPlayerKills struct {
	AlbionService *services.AlbionService
}

func (uc *FetchPlayerKills) Handle(playerId string) ([]entities.PlayerKill, error) {
	return uc.AlbionService.FetchPlayerKills(playerId)
}
