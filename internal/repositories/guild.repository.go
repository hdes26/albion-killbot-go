package repositories

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
)

type GuildRepository struct {
	AlbionService *services.AlbionService
}

func (r *GuildRepository) GetMembers() ([]entities.GuildMember, error) {
	// Aquí se podría acceder a la API o base de datos para obtener los miembros
	// Ejemplo de simulación:
	members := []entities.GuildMember{
		{Id: "1", Name: "Player1"},
		{Id: "2", Name: "Player2"},
	}
	return members, nil
}
