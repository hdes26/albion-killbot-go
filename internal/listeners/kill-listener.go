package listeners

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/usecases"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type KillListener struct {
	FetchPlayerKills  usecases.FetchPlayerKills
	FetchGuildMembers usecases.FetchGuildMembers
}

func (l *KillListener) Start(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	log.Println("Listener started, waiting for events...")

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

	// Obtener miembros del gremio con reintentos
	var members []entities.GuildMember
	var err error
	for {
		members, err = l.FetchGuildMembers.Handle()
		if err == nil {
			break // Si no hay error, salimos del bucle
		}
		log.Printf("Error al obtener miembros del gremio: %v. Intentando de nuevo en 1 minuto...", err)
		time.Sleep(1 * time.Minute) // Esperar 1 minuto antes de intentar nuevamente
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

	// Manejar los resultados
	for result := range results {
		// Aquí puedes agregar la lógica para enviar mensajes o cualquier otra acción
		log.Printf("Received result: %v", result)
	}

	// Manejar los errores
	for err := range errors {
		log.Printf("Error processing kill: %v", err)
	}
}

// worker procesa las kills de cada miembro en paralelo
func (l *KillListener) worker(
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
			return
		case memberId, ok := <-tasks:
			if !ok {
				return
			}

			kills, err := l.FetchPlayerKills.Handle(memberId)
			if err != nil {
				errors <- fmt.Errorf("error fetching kills for player %s: %w", memberId, err)
				continue
			}
			results <- kills
		}
	}
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
