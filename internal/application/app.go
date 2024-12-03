package application

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
	"albion-killbot/internal/usecases"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type App struct {
	fetchGuildMembers usecases.FetchGuildMembers
	fetchPlayerKills  usecases.FetchPlayerKills
	previousKills     map[string]entities.PlayerKill
}

func NewApp(albionService services.AlbionService) *App {
	return &App{
		fetchGuildMembers: usecases.FetchGuildMembers{AlbionService: &albionService},
		fetchPlayerKills:  usecases.FetchPlayerKills{AlbionService: &albionService},
		previousKills:     make(map[string]entities.PlayerKill), // Inicializar el mapa para las kills previas
	}
}

func (app *App) Run(ctx context.Context) error {
	// Crear un ticker para ejecutar la verificación periódica
	ticker := time.NewTicker(15 * time.Minute) // Ejecutar cada 15 minutos
	defer ticker.Stop()

	log.Println("Aplicación iniciada, esperando eventos del ticker...")

	// Ejecutar la lógica continuamente mientras el ticker esté activo
	for {
		select {
		case <-ctx.Done():
			// Si el contexto se cancela, salimos del bucle
			log.Println("Proceso detenido")
			return nil
		case <-ticker.C:
			// Ejecutar la lógica cada vez que el ticker genere un evento
			log.Println("Ticker disparado, ejecutando lógica")
			if err := app.fetchAndProcessKills(ctx); err != nil {
				log.Printf("Error procesando kills: %v", err)
			}
		}
	}
}

func (app *App) fetchAndProcessKills(ctx context.Context) error {
	// Obtener miembros del gremio con reintentos
	var members []entities.GuildMember
	var err error
	for {
		members, err = app.fetchGuildMembers.Handle()
		if err == nil {
			break // Si no hay error, salimos del bucle
		}
		log.Printf("Error al obtener miembros del gremio: %v. Intentando de nuevo en 1 minuto...", err)
		time.Sleep(1 * time.Minute) // Esperar 1 minuto antes de intentar nuevamente
	}

	// Canales para manejar tareas y resultados
	tasks := make(chan string)
	results := make(chan []entities.PlayerKill)
	errors := make(chan error)

	// Enviar IDs de miembros al canal de tareas
	go func() {
		for _, member := range members {
			tasks <- member.Id
		}
		close(tasks) // Cerrar el canal solo después de enviar todas las tareas
	}()

	// Número de trabajadores concurrentes
	const maxWorkerCount = 20
	workerCount := calculateWorkers(len(members))
	workerCount = min(workerCount, maxWorkerCount)

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go app.worker(ctx, tasks, results, errors, &wg)
	}

	// Manejar los resultados de los workers
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Manejar los resultados y errores
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	go func() {
		for err := range errors {
			log.Printf("Error procesando kills: %v", err)
		}
	}()

	return nil
}

func (app *App) worker(
	ctx context.Context,
	tasks <-chan string,
	results chan<- []entities.PlayerKill,
	errors chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return // Contexto cancelado o tiempo excedido
		case memberId, ok := <-tasks:
			if !ok {
				return // No más tareas
			}
			// Procesar kills
			kills, err := app.fetchPlayerKills.Handle(memberId)
			if err != nil {
				errors <- fmt.Errorf("error en miembro %s: %w", memberId, err)
				continue
			}
			if len(kills) > 0 && kills[0].EventId != app.previousKills[memberId].EventId {
				// Si la más reciente kill no es igual a la anterior, es nueva
				// Convertir EventId a string para poder imprimirlo correctamente
				eventId := strconv.Itoa(kills[0].EventId)
				fmt.Printf("¡Nueva kill detectada para %s EventId: %s\n", memberId, eventId)

				// Actualizar la kill previa
				app.previousKills[memberId] = kills[0]
			}
			results <- kills
		}
	}
}

func calculateWorkers(taskQueueSize int) int {
	switch {
	case taskQueueSize > 100:
		return 20
	case taskQueueSize > 50:
		return 15
	case taskQueueSize > 0:
		return 10
	default:
		return 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
