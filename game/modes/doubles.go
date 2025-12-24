package modes

type Doubles struct {
}

func (d Doubles) String() string {
	return "Doubles"
}

func (d Doubles) MinimumTotalPlayers() int {
	return 2
}

func (d Doubles) MaximumTotalPlayers() int {
	return 16
}

func (d Doubles) NumberOfPlayersPerTeam() int {
	return 2
}
