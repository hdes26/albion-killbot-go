package entities

// Event representa un evento de combate en el juego
type Event struct {
	EventID     string
	KillerGuild string
	KillerName  string
	VictimGuild string
	VictimName  string
}
