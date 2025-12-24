package modes

var registry = map[string]Mode{
	Normal{}.String(): Normal{},
}

func RegisterMode(mode Mode) {
	registry[mode.String()] = mode
}

func GetModeFromString(mode string) Mode {
	if m, ok := registry[mode]; ok {
		return m
	}
	return nil
}
