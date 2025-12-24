package modes

type Duels struct {
}

func (d Duels) String() string {
	return "Duels"
}

func (d Duels) MinimumTotalPlayers() int {
	return 2
}

func (d Duels) MaximumTotalPlayers() int {
	return 2
}

func (d Duels) NumberOfPlayersPerTeam() int {
	return 1
}
