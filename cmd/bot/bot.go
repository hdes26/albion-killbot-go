package main

import (
	"albion-killbot/internal/application"
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

	fmt.Println("Bot running")

	// Crear el bot
	bot := application.NewBot(botToken)
	if bot == nil {
		log.Fatalf("Error al inicializar el bot")
	}

	// Crear un contexto sin timeout (la aplicación siempre estará ejecutándose)
	ctx := context.Background()

	// Ejecutar la lógica de la aplicación

	// Ejecutar el bot, en caso de que la aplicación haya finalizado sin errores
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("Error en la ejecución del bot: %v", err)
	}

	// Si ambos procesos se ejecutan correctamente
	fmt.Println("Bot ejecutado con éxito")
}
