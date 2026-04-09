package command

import (
	"github.com/blockbrawn/game"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PauseCommand struct {
}

func (pc PauseCommand) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	series := game.GetGame().StateSeries

	if series.IsPaused() {
		output.Error("El juego ya está pausado")

		return
	}

	series.SetPaused(true)
	output.Print(text.Colourf("<green>El juego se ha pausado exitosamente</green>"))
}
