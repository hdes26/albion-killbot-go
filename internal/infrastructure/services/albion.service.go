package services

import (
	"albion-killbot/internal/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AlbionAPIService es un servicio centralizado para interactuar con la API de Albion
type AlbionService struct{}

type GuildResponse struct {
	Guilds  []entities.Guild `json:"guilds"`
	Players []interface{}    `json:"players"` // Si no necesitas los players, puedes ignorarlos.
}

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

func (service *AlbionService) FetchGuildByName(guildName string) (entities.Guild, error) {
	const maxRetries = 5
	const retryInterval = 30 * time.Second

	var guilds []entities.Guild
	var lastError error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Intentando obtener guilds: intento %d/%d\n", attempt, maxRetries)

		req, err := http.NewRequest("GET", "https://gameinfo.albiononline.com/api/gameinfo/search", nil)
		if err != nil {
			fmt.Println("Error", err)
		}
		q := req.URL.Query()
		q.Add("q", guildName)
		req.URL.RawQuery = q.Encode()

		resp, err := http.Get(req.URL.String())
		if err != nil {
			lastError = err
			fmt.Printf("Error en la solicitud: %s. Reintentando...\n", err)
			time.Sleep(retryInterval)
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastError = err
			fmt.Printf("Error leyendo el cuerpo de la respuesta: %s. Reintentando...\n", err)
			time.Sleep(retryInterval)
			continue
		}

		// Verifica si la respuesta parece HTML (por ejemplo, una página de error)
		if len(body) > 0 && body[0] == '<' {
			fmt.Println("Parece que la respuesta es HTML, no JSON.")
			time.Sleep(retryInterval)
			continue
		}

		var response GuildResponse
		if err := json.Unmarshal(body, &response); err != nil {
			lastError = err
			fmt.Printf("Error deserializando la respuesta JSON: %s. Reintentando...\n", err)
			time.Sleep(retryInterval)
			continue
		}

		// Retorna los guilds si hay resultados
		if len(response.Guilds) > 0 {
			guilds = response.Guilds
			return guilds[0], nil
		}

		// Si no se encuentran guilds, imprime mensaje y reintenta
		fmt.Println("No se encontraron guilds en la respuesta. Reintentando...")
		time.Sleep(retryInterval)
	}

	// Si todos los intentos fallan, retorna el último error
	if lastError == nil {
		lastError = errors.New("no se encontraron guilds después de varios intentos")
	}
	return entities.Guild{}, lastError
}
