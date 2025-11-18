package skins

import (
	"log"

	"github.com/BlockBrawn/game/game/utils/randomskins/skin"
)

var SkinManager skin.Manager

func init() {
	skinManager, err := skin.SetupSkinManager()
	if err != nil {
		log.Fatalf("Could not set up a new skin manager: %v", err)
	}

	SkinManager = skinManager
}
