package modes

var registry = map[string]Mode{
	Normal{}.ID(): Normal{},
}

func RegisterMode(mode Mode) {
	registry[mode.ID()] = mode
}

func GetModeFromString(mode string) Mode {
	if m, ok := registry[mode]; ok {
		return m
	}
	return nil
}
