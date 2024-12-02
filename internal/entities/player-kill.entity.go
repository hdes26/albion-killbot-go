package entities

type PlayerKill struct {
	GroupMemberCount     int     `json:"groupMemberCount"`
	NumberOfParticipants int     `json:"numberOfParticipants"`
	EventId              int     `json:"EventId"`
	TimeStamp            string  `json:"TimeStamp"`
	Version              int     `json:"Version"`
	Killer               Player  `json:"Killer"`
	Victim               Player  `json:"Victim"`
	TotalVictimKillFame  int     `json:"TotalVictimKillFame"`
	Location             *string `json:"Location"`
	GvGMatch             *string `json:"GvGMatch"`
	BattleId             int     `json:"BattleId"`
	KillArea             string  `json:"KillArea"`
	Category             *string `json:"Category"`
	Type                 string  `json:"Type"`
}
