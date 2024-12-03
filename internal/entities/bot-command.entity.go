package entities

// CommandOption representa las opciones de un comando
type CommandOption struct {
	Type        int      `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Choices     []Choice `json:"choices,omitempty"`
}

// Choice representa las opciones dentro de un comando
type Choice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Command representa un comando completo de Discord
type Command struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     []CommandOption `json:"options"`
}

var ChannelTypes = []Choice{
	{Name: "Kill", Value: "kill"},
	{Name: "Death", Value: "death"},
}

// Commands contiene todos los comandos del bot
var Commands = []Command{
	{
		Name:        "killboard",
		Description: "Install Guild Kill Registration",
		Options: []CommandOption{
			{
				Type:        3,
				Name:        "guild",
				Description: "Guild name",
				Required:    true,
			},
		},
	},
	{
		Name:        "set",
		Description: "Install Channel Kill Registration",
		Options: []CommandOption{
			{
				Type:        3,
				Name:        "type",
				Description: "Channel type",
				Required:    true,
				Choices:     ChannelTypes,
			},
		},
	},
}
