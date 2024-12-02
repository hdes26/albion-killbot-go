package entities

type Player struct {
	AverageItemPower float64   `json:"AverageItemPower"` // Poder promedio del equipo del jugador
	Equipment        Equipment `json:"Equipment"`        // Equipamiento del jugador
	Inventory        []Item    `json:"Inventory"`        // Inventario del jugador
	Name             string    `json:"Name"`             // Nombre del jugador
	Id               string    `json:"Id"`
	GuildName        string    `json:"GuildName"`    // Nombre de la hermandad
	GuildId          string    `json:"GuildId"`      // ID de la hermandad
	AllianceName     string    `json:"AllianceName"` // Nombre de la alianza
	AllianceId       string    `json:"AllianceId"`   // ID de la alianza
	AllianceTag      string    `json:"AllianceTag"`  // Etiqueta de la alianza
	Avatar           string    `json:"Avatar"`       // Avatar del jugador
	AvatarRing       string    `json:"AvatarRing"`   // Anillo del avatar
	DeathFame        int       `json:"DeathFame"`    // Fama por muertes
	KillFame         int       `json:"KillFame"`     // Fama por kills
	FameRatio        float64   `json:"FameRatio"`    // Ratio de fama
}

type Equipment struct {
	MainHand Item  `json:"MainHand"`
	OffHand  *Item `json:"OffHand"`
	Head     Item  `json:"Head"`
	Armor    Item  `json:"Armor"`
	Shoes    Item  `json:"Shoes"`
	Bag      Item  `json:"Bag"`
	Cape     Item  `json:"Cape"`
	Mount    Item  `json:"Mount"`
	Potion   Item  `json:"Potion"`
	Food     Item  `json:"Food"`
}

type Item struct {
	Type          string      `json:"Type"`          // Tipo del equipo
	Count         int         `json:"Count"`         // Cantidad
	Quality       int         `json:"Quality"`       // Calidad del equipo
	ActiveSpells  []string    `json:"ActiveSpells"`  // Hechizos activos
	PassiveSpells []string    `json:"PassiveSpells"` // Hechizos pasivos
	LegendarySoul interface{} `json:"LegendarySoul"` // Alma legendaria (ahora acepta cualquier tipo)
}
