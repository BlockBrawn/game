package command

import (
	"github.com/blockbrawn/game"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SkipCommand struct {
}

func (sc SkipCommand) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	game.GetGame().StateSeries.Skip()
	output.Print(text.Colourf("<green>El juego ha saltado al siguiente estado exitosamente</green>"))
}
