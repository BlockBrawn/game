package modes

var registry = map[string]Mode{
	Normal{}.String(): Normal{},
	Solo{}.String():   Solo{},
	Duo{}.String():    Duo{},
	Squad{}.String():  Squad{},
	FFA{}.String():    FFA{},
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
