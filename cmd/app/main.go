package main

import (
	"albion-killbot/internal/application"
	"albion-killbot/internal/infrastructure/services"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("El token del bot no está definido en el archivo .env")
	}

	fmt.Println("App running")

	// Inicializar servicios
	albionService := services.NewAlbionService()

	// Crear el servicio de aplicación que orquesta los casos de uso
	app := application.NewApp(*albionService)

	// Crear el bot
	bot := application.NewBot(botToken)
	if bot == nil {
		log.Fatalf("Error al inicializar el bot")
	}

	// Crear un contexto sin timeout (la aplicación siempre estará ejecutándose)
	ctx := context.Background()

	// Ejecutar la lógica de la aplicación
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Error en la ejecución de la aplicación: %v", err)
	}

	// Ejecutar el bot, en caso de que la aplicación haya finalizado sin errores
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("Error en la ejecución del bot: %v", err)
	}

	// Si ambos procesos se ejecutan correctamente
	fmt.Println("Aplicación ejecutada con éxito")
}
