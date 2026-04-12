package bots

import (
	"fmt"

	"github.com/blockbrawn/game"
	"github.com/blockbrawn/game/participant"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type FillSubCommand struct {
	Fill cmd.SubCommand `cmd:"fill"`
}

func (asc FillSubCommand) Run(source cmd.Source, _ *cmd.Output, _ *world.Tx) {
	p := source.(*player.Player)

	g := game.GetGame()
	needed := g.Settings.Mode.MaximumTotalPlayers() - g.GetParticipantLenByState(participant.StateAlive)

	p.ExecuteCommand(fmt.Sprintf("/bots add %d", needed))
}
