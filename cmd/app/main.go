package main

import (
	"albion-killbot/internal/application"
	database "albion-killbot/internal/infrastructure/db"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("El token del bot no está definido en el archivo .env")
	}

	// Conectar a MongoDB
	mongoUri := os.Getenv("MG_URI")
	dbClient, err := database.Connect(mongoUri)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer dbClient.Disconnect()

	fmt.Println("App running")

	// Crear la aplicación y pasar la conexión de la base de datos
	app := application.NewApp(botToken, dbClient) // Pasa la conexión de DB

	// Crear un contexto sin timeout
	ctx := context.Background()

	// Ejecutar la lógica de la aplicación
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Error en la ejecución de la aplicación: %v", err)
	}

	fmt.Println("Aplicación ejecutada con éxito")
}
