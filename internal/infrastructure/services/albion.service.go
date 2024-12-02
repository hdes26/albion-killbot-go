package services

import (
	"albion-killbot/internal/entities"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AlbionAPIService es un servicio centralizado para interactuar con la API de Albion
type AlbionService struct{}

// NewAlbionAPIService crea una nueva instancia del servicio de API de Albion
func NewAlbionService() *AlbionService {
	return &AlbionService{}
}

func (service *AlbionService) FetchGuildMembers() ([]entities.GuildMember, error) {
	const guildEndpoint = "https://gameinfo.albiononline.com/api/gameinfo/guilds/Amt2FdNMRTWRYmhaYnibvQ/members"
	resp, err := http.Get(guildEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var members []entities.GuildMember
	if err := json.Unmarshal(body, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func (service *AlbionService) FetchPlayerKills(playerId string) ([]entities.PlayerKill, error) {
	killEndpoint := fmt.Sprintf("https://gameinfo.albiononline.com/api/gameinfo/players/%s/kills", playerId)
	resp, err := http.Get(killEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kills []entities.PlayerKill
	if err := json.Unmarshal(body, &kills); err != nil {
		return nil, err
	}

	return kills, nil
}
