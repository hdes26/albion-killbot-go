package usecases

import (
	"albion-killbot/internal/infrastructure/services"
	"fmt"
)

// SendInteractionUseCase es el caso de uso para enviar interacciones al servidor de Discord
type SendKillEventUseCase struct {
	BotService *services.BotService // Inyectamos el servicio del bot
}

// NewSendInteractionUseCase crea una nueva instancia del caso de uso
func NewSendKillEeventUseCase(botService *services.BotService) *SendKillEventUseCase {
	return &SendKillEventUseCase{
		BotService: botService,
	}
}

// Handle envía la interacción al servidor de Discord
func (uc *SendKillEventUseCase) Handle(channelId string, message string) error {
	fmt.Println(channelId, message)

	fmt.Println("--------------------")

	err := uc.BotService.SendInteractionToServer(channelId, message)
	if err != nil {
		return fmt.Errorf("error al enviar interacción: %v", err)
	}
	return nil
}
