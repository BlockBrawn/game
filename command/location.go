package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type LocationCommand struct{}

func (lc LocationCommand) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	p := source.(*player.Player)

	pos := p.Position()
	rot := p.Rotation()

	output.Print(text.Colourf(
		"<aqua>Posición: </aqua><grey>%.2f, %.2f, %.2f</grey>\n<aqua>Rotación: </aqua><grey>%.2f, %.2f</grey>",
		pos.X(), pos.Y(), pos.Z(), rot.Yaw(), rot.Pitch(),
	))
}
