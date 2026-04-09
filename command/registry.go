package command

import (
	"github.com/blockbrawn/game/command/bots"
	"github.com/df-mc/dragonfly/server/cmd"
)

func RegisterDevCommands() {
	cmd.Register(cmd.New(
		"pause",
		"Pausar tu juego",
		nil,
		PauseCommand{},
	))
	cmd.Register(cmd.New(
		"resume",
		"Reanudar tu juego",
		nil,
		ResumeCommand{},
	))
	cmd.Register(cmd.New(
		"skip",
		"Saltar al siguiente estado",
		nil,
		SkipCommand{},
	))
	cmd.Register(cmd.New(
		"location",
		"Mostrar tu posición",
		[]string{"loc"},
		LocationCommand{},
	))
	cmd.Register(cmd.New(
		"bots",
		"Administrar bots del juego",
		nil,
		bots.AddSubCommand{},
		bots.RemoveSubCommand{},
		bots.FillSubCommand{},
	))
}
