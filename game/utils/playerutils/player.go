package playerutils

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
)

type ResetOpts struct {
	ClearEffects      bool
	Extinguish        bool
	HealAmount        float64
	MaxHealth         float64
	ResetArmour       bool
	ResetFallDistance bool
	ResetFood         bool
	ResetHealth       bool
	ResetInventory    bool
	ResetMobility     bool
	ResetScale        bool
	ResetXP           bool
}

func DefaultResetOpts() ResetOpts {
	return ResetOpts{
		ClearEffects:      true,
		Extinguish:        true,
		HealAmount:        20,
		MaxHealth:         20,
		ResetArmour:       true,
		ResetFallDistance: true,
		ResetFood:         true,
		ResetHealth:       true,
		ResetInventory:    true,
		ResetMobility:     true,
		ResetScale:        true,
		ResetXP:           true,
	}
}

type ResetPlayerHealSource struct{}

func (ResetPlayerHealSource) HealingSource() {}

func ResetPlayer(p *player.Player, opts *ResetOpts) {
	if opts == nil {
		def := DefaultResetOpts()
		opts = &def
	}

	if opts.ResetInventory {
		p.Inventory().Clear()
	}
	if opts.ResetArmour {
		p.Armour().Clear()
	}
	p.MoveItemsToInventory()
	p.SetHeldItems(item.Stack{}, item.Stack{})

	if opts.Extinguish {
		p.Extinguish()
	}
	if opts.ClearEffects {
		for _, effect := range p.Effects() {
			p.RemoveEffect(effect.Type())
		}
	}
	if opts.ResetFallDistance {
		p.ResetFallDistance()
	}
	if opts.ResetXP {
		p.SetExperienceProgress(0)
		p.SetExperienceLevel(0)
	}
	if opts.ResetHealth {
		p.SetMaxHealth(opts.MaxHealth)
		p.Heal(opts.HealAmount, ResetPlayerHealSource{})
	}
	if opts.ResetScale {
		p.SetScale(1)
	}
	if opts.ResetFood {
		p.SetFood(20)
	}
	if opts.ResetMobility {
		p.SetMobile()
	}
}

func SetRotation(p *player.Player, yaw, pitch float64) {
	currentRotation := p.Rotation()

	deltaYaw := yaw - currentRotation.Yaw()
	deltaPitch := pitch - currentRotation.Pitch()

	p.Move(mgl64.Vec3{}, deltaYaw, deltaPitch)
	p.Teleport(p.Position())
}
