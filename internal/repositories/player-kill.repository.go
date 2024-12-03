package repositories

import "albion-killbot/internal/entities"

type MemberKillRepository interface {
	FetchMemberKills(playerId string) ([]entities.PlayerKill, error)
}
