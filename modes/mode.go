package modes

type Mode interface {
	ID() string
	Name() string
	MinimumTotalPlayers() int
	MaximumTotalPlayers() int
	NumberOfPlayersPerTeam() int
}
