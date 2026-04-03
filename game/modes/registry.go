package modes

import "github.com/blockbrawn/game/game/utils/maputils"

type key struct {
	gameID string
	modeID string
}

var registry = maputils.NewMap[key, Mode]()

func RegisterMode(gameID string, mode Mode) {
	k := key{
		gameID: gameID,
		modeID: mode.ID(),
	}
	registry.Store(k, mode)
}

func GetMode(gameID, modeID string) Mode {
	if m, ok := registry.Load(key{
		gameID: gameID,
		modeID: modeID,
	}); ok {
		return m
	}
	return nil
}
