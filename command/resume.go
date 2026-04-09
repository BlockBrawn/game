package command

import (
	"github.com/blockbrawn/game"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ResumeCommand struct {
}

func (rc ResumeCommand) Run(_ cmd.Source, output *cmd.Output, _ *world.Tx) {
	series := game.GetGame().StateSeries

	if !series.IsPaused() {
		output.Error("El juego no está pausado")

		return
	}

	series.SetPaused(false)
	output.Print(text.Colourf("<green>El juego se ha reanudado exitosamente</green>"))
}
