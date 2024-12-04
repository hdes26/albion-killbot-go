package listeners

import (
	"albion-killbot/internal/entities"
	dbrepositories "albion-killbot/internal/infrastructure/db/repositories"
	"albion-killbot/internal/infrastructure/services"
	"albion-killbot/internal/usecases"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type KillListener struct {
	FetchPlayerKills  usecases.FetchPlayerKills
	FetchGuildMembers usecases.FetchGuildMembers
	SendKillEvent     usecases.SendKillEventUseCase
	ChannelRepo       *dbrepositories.ChannelRepository
	Botservice        *services.BotService
}

func NewKillListener(session *discordgo.Session, channelRepo *dbrepositories.ChannelRepository) *KillListener {

	// Inicializar los casos de uso
	fetchPlayerKills := usecases.FetchPlayerKills{}
	fetchGuildMembers := usecases.FetchGuildMembers{}
	botservice := services.BotService{
		Session: session,
	}

	sendKillEvent := usecases.SendKillEventUseCase{
		BotService: &botservice,
	}

	// Crear y devolver el KillListener
	return &KillListener{
		FetchPlayerKills:  fetchPlayerKills,
		FetchGuildMembers: fetchGuildMembers,
		SendKillEvent:     sendKillEvent,
		ChannelRepo:       channelRepo,
		Botservice:        &botservice,
	}
}
func (l *KillListener) Start(ctx context.Context, channelRepo *dbrepositories.ChannelRepository) {

	l.processKills(ctx)

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	log.Println("Listener started, waiting for events...")
	l.processKills(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("Listener stopped")
			return
		case <-ticker.C:
			log.Println("Ticker triggered, processing kills...")
			l.processKills(ctx)
		}
	}
}

func (l *KillListener) processKills(ctx context.Context) {
	var err error
	channels, err := l.ChannelRepo.FindChannels()
	if err != nil {
		log.Printf("Error buscando canales: %v", err)
		return
	}
	for _, channel := range channels {

		members, err := l.fetchGuildMembersWithRetry(channel.Guild.Id)
		if err != nil || len(members) == 0 {
			log.Printf("No se encontraron miembros en el gremio %s (Canal: %s). Pasando al siguiente canal...", channel.Guild.Name, *channel.Channel)
			continue
		}

		tasks := make(chan string)
		results := make(chan []entities.PlayerKill)
		errors := make(chan error)

		go func() {
			for _, member := range members {
				tasks <- member.Id
			}
			close(tasks)
		}()

		// Calcular el número de workers de manera dinámica
		workerCount := calculateWorkers(len(members))
		var wg sync.WaitGroup

		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go l.worker(ctx, tasks, results, errors, &wg)
		}

		go func() {
			wg.Wait()
			close(results)
			close(errors)
		}()

		l.handleResultsAndErrors(channel.ChannelID, results, errors)
	}
	fmt.Print("Sending kills...")
}

func (l *KillListener) fetchGuildMembersWithRetry(guildID string) ([]entities.GuildMember, error) {
	var members []entities.GuildMember
	var err error

	// Reintentos con un límite de 5
	for i := 0; i < 5; i++ {
		members, err = l.FetchGuildMembers.Handle(guildID)
		if err == nil {
			return members, nil
		}
		log.Printf("Error al obtener miembros del gremio: %v. Intentando nuevamente en 1 minuto...", err)
		time.Sleep(1 * time.Minute)
	}
	return nil, fmt.Errorf("error al obtener miembros después de 5 intentos: %w", err)
}

// worker procesa las kills de cada miembro en paralelo
func (l *KillListener) worker(
	ctx context.Context,
	tasks <-chan string,
	results chan<- []entities.PlayerKill,
	errors chan<- error,
	wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done(): // Si el contexto se cancela, terminamos el worker
			log.Println("Worker cancelado por contexto")
			return
		case memberId, ok := <-tasks:
			if !ok {
				// Si el canal de tareas se cierra, terminamos el worker
				log.Println("Canal de tareas cerrado, finalizando worker")
				return
			}

			// Obtener las kills del jugador
			kills, err := l.FetchPlayerKills.Handle(memberId)
			if err != nil {
				// Si ocurre un error, enviamos el error y continuamos con el siguiente miembro
				errors <- fmt.Errorf("error fetching kills for player %s: %w", memberId, err)
				continue // No enviamos resultados si hay error
			}

			// Enviar los resultados obtenidos
			select {
			case results <- kills:
			case <-ctx.Done():
				// Si el contexto se cancela mientras enviamos, terminamos la operación
				log.Println("Worker cancelado durante el envío de resultados")
				return
			}
		}
	}
}

func (l *KillListener) handleResultsAndErrors(channelId string, results <-chan []entities.PlayerKill, errors <-chan error) error {
	for playersKills := range results {
		for _, kill := range playersKills {

			err := l.SendKillEvent.Handle(channelId, kill)

			if err != nil {
				log.Printf("Error al enviar evento de kill: %v", err)
				return err // Propagar el error a la capa superior
			}
		}
	}

	for err := range errors {
		log.Printf("Error procesando kill: %v", err)
		return err
	}

	return nil
}

// calculateWorkers determina la cantidad de workers según el tamaño de la cola de tareas
func calculateWorkers(taskQueueSize int) int {
	// Define el número máximo de workers
	const maxWorkers = 20

	// Calcula el número de workers según el tamaño de la cola de tareas
	var workers int
	switch {
	case taskQueueSize > 100:
		workers = 20
	case taskQueueSize > 50:
		workers = 15
	case taskQueueSize > 0:
		workers = 10
	default:
		workers = 1
	}

	// Asegura que no se superen los workers máximos
	return min(workers, maxWorkers)
}

// min devuelve el valor menor entre dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
