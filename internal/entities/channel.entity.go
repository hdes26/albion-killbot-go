package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Guild representa un guild en la base de datos.
type Channel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`  // El ID de MongoDB para el guild.
	ChannelID string             `bson:"channel_id"`     // ID del canal donde se configuró el guild.
	Channel   string             `bson:"channel"`        // Nombre del canal donde se configuró el guild.
	Guild     *Guild             `json:"Guild"`          // Equipamiento del jugador
	Type      *string            `bson:"type,omitempty"` // Tipo del canal, como "killboard", "event", etc.
	Active    *bool              `bson:"active"`         // Indica si el guild está activo.
	CreatedAt time.Time          `bson:"created_at"`     // Fecha de creación del guild.
	UpdatedAt time.Time          `bson:"updated_at"`     // Fecha de la última actualización.
}

type Guild struct {
	Id           string `json:"Id"`
	Name         string `json:"Name"`
	AllianceId   string `json:"AllianceId"`
	AllianceName string `json:"AllianceName"`
	KillFame     int    `json:"KillFame"`
	DeathFame    int    `json:"DeathFame"`
}
