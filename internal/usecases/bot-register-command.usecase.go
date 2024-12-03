package usecases

import (
	"albion-killbot/internal/entities"
	"albion-killbot/internal/infrastructure/services"
	"fmt"
)

// BotRegisterCommandsUseCase maneja los casos de uso relacionados con el registro de comandos del bot
type BotRegisterCommandsUseCase struct {
	BotService *services.BotService
}

// RegisterCommands maneja la validación y el registro de comandos
func (uc *BotRegisterCommandsUseCase) Handle() error {
	// Validación de comandos
	if len(entities.Commands) == 0 {
		return fmt.Errorf("no hay comandos definidos para registrar")
	}

	// Validar cada comando
	for _, cmd := range entities.Commands {
		if err := uc.validateCommand(cmd); err != nil {
			return err
		}
	}

	// Delegar al servicio para registrar los comandos en Discord
	for _, cmd := range entities.Commands {
		err := uc.BotService.RegisterCommand(cmd) // Usar la sesión desde el campo de la estructura
		if err != nil {
			return fmt.Errorf("error al registrar el comando '%s': %v", cmd.Name, err)
		}
	}

	return nil
}

// validateCommand valida que el comando tenga un nombre y descripción válidos
func (uc *BotRegisterCommandsUseCase) validateCommand(cmd entities.Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("el comando '%s' debe tener un nombre", cmd.Name)
	}
	if cmd.Description == "" {
		return fmt.Errorf("el comando '%s' debe tener una descripción", cmd.Name)
	}

	// Puedes agregar más validaciones aquí si es necesario

	return nil
}
