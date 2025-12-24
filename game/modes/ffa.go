package modes

type FFA struct {
}

func (f FFA) String() string {
	return "ffa"
}

func (f FFA) MinimumTotalPlayers() int {
	return 0
}

func (f FFA) MaximumTotalPlayers() int {
	return 100
}

func (f FFA) NumberOfPlayersPerTeam() int {
	return 1
}
