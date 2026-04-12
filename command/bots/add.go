package bots

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/blockbrawn/game"
	"github.com/blockbrawn/game/handler_custom"
	"github.com/blockbrawn/game/mechanic/bot"
	"github.com/blockbrawn/game/participant"
	"github.com/blockbrawn/game/skins"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/npc"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type AddSubCommand struct {
	Add    cmd.SubCommand    `cmd:"add"`
	Number cmd.Optional[int] `cmd:"number"`
}

func (asc AddSubCommand) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	p := source.(*player.Player)

	number, ok := asc.Number.Load()
	if number <= 0 || !ok {
		number = 1
	}

	g := game.GetGame()
	needed := g.Settings.Mode.MaximumTotalPlayers() - g.GetParticipantLenByState(participant.StateAlive)

	if needed <= 0 {
		output.Error("No se necesitan bots, el juego ya está lleno.")
		return
	}

	if number > needed {
		output.Print(text.Colourf("<yellow>Solo puedes agregar %d bot(s) en este momento.</yellow>", needed))
		number = needed
	}

	go func() {
		for i := 0; i < number; i++ {
			if err := skins.SkinManager.GenerateSkin(i); err != nil {
				log.Fatalf("Could not generate a new skin: %v", err)
			}
			p.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				skin := npc.MustSkin(npc.MustParseTexture(path.Join(".", "skins", fmt.Sprintf("skin_%v.png", i))), npc.DefaultModel)
				bot.AddBot(tx, p.Position(), p.Rotation(), skin, func(b *player.Player) {
					b.Handle(g.PlayerHandler)
					b.Handler().(handler_custom.JoinHandler).HandleJoin(b)
				})
			})
			time.Sleep(500 * time.Millisecond)
		}
	}()

	output.Print(text.Colourf("<green>%d bot(s) agregado(s) exitosamente</green>", number))
}
