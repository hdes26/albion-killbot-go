package main

import (
	"albion-killbot/internal/application"
	"albion-killbot/internal/infrastructure/services"
	"context"
	"fmt"
	"log"
)

func main() {
	fmt.Println("app running")

	// Inicializar servicios
	albionService := services.NewAlbionService()

	// Crear el servicio de aplicación que orquesta los casos de uso
	app := application.NewApp(*albionService)

	// Crear un contexto sin timeout (la aplicación siempre estará ejecutándose)
	ctx := context.Background()

	// Ejecutar la lógica de la aplicación
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Error en la ejecución de la aplicación: %v", err)
	}

	fmt.Println("Aplicación ejecutada con éxito")
}
