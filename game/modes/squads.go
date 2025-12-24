package modes

type Squads struct {
}

func (s Squads) String() string {
	return "Squads"
}

func (s Squads) MinimumTotalPlayers() int {
	return 2
}

func (s Squads) MaximumTotalPlayers() int {
	return 16
}

func (s Squads) NumberOfPlayersPerTeam() int {
	return 4
}
