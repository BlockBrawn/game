package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/blockbrawn/game/config"
	"github.com/blockbrawn/game/handler_custom"
	"github.com/blockbrawn/game/participant"
	"github.com/blockbrawn/game/settings"
	"github.com/blockbrawn/game/team"
	"github.com/blockbrawn/game/utils/maputils"
	"github.com/blockbrawn/game/utils/ziputils"
	"github.com/df-mc/dragonfly/server/item/inventory"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/josscoder/fsmgo/state"
)

var gameInstance *Game

type Game struct {
	id       uuid.UUID
	Settings *settings.Settings
	Teams    []*team.Team

	PlayerHandler    handler_custom.JoinHandler
	InventoryHandler inventory.Handler

	StateSeries  *state.ScheduledStateSeries
	Participants *maputils.Map[uuid.UUID, *participant.Participant]

	MapLoaded bool
	mapConfig config.MapData

	World       *world.World
	WorldFolder string
}

func NewGame(settings *settings.Settings, teams []*team.Team, states []state.State, playerHandler handler_custom.JoinHandler, invHandler inventory.Handler) *Game {
	if playerHandler == nil {
		panic("player handler cannot be nil")
	}

	game := &Game{
		id:               uuid.New(),
		Settings:         settings,
		Teams:            teams,
		PlayerHandler:    playerHandler,
		InventoryHandler: invHandler,
		StateSeries:      state.NewScheduledStateSeries(states),
		Participants:     maputils.NewMap[uuid.UUID, *participant.Participant](),
	}
	gameInstance = game
	return game
}

func GetGame() *Game {
	return gameInstance
}

func (g *Game) Start() error {
	g.StateSeries.Start()
	return nil
}

func (g *Game) LoadGameMapWithConfig(config config.MapData) error {
	worldsDir := path.Join(".", "worlds")
	if entries, err := os.ReadDir(worldsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				oldWorldPath := filepath.Join(worldsDir, entry.Name())
				if err := os.RemoveAll(oldWorldPath); err != nil {
					fmt.Println("warning: failed to remove old world folder:", oldWorldPath, err)
				}
			}
		}
	}

	if g.MapLoaded {
		return errors.New("map already loaded")
	}

	mapDir := filepath.Join(g.Settings.MapsFolder, g.Settings.MapName)
	worldZip := filepath.Join(mapDir, "world.zip")

	stat, err := os.Stat(worldZip)
	if err != nil || stat.IsDir() {
		return fmt.Errorf("world.zip not found or is not a file: %s", worldZip)
	}

	rawConfig, err := os.ReadFile(filepath.Join(mapDir, "config.json"))
	if err != nil {
		return fmt.Errorf("failed reading config.json: %w", err)
	}

	if err := json.Unmarshal(rawConfig, config); err != nil {
		return fmt.Errorf("failed unmarshalling config.json: %w", err)
	}

	g.WorldFolder = path.Join(".", "worlds", g.id.String())

	if err := ziputils.UnZipFile(worldZip, g.WorldFolder); err != nil {
		return fmt.Errorf("failed to copy world: %w", err)
	}

	g.MapLoaded = true
	g.mapConfig = config

	return nil
}

func GetMapData[TM config.MapData]() TM {
	if !gameInstance.MapLoaded {
		panic("map is not loaded")
	}
	return gameInstance.mapConfig.(TM)
}

func (g *Game) GetParticipant(p *player.Player) *participant.Participant {
	pt, _ := g.Participants.Load(p.UUID())
	return pt
}

func (g *Game) GetParticipants() iter.Seq[*participant.Participant] {
	return func(yield func(*participant.Participant) bool) {
		for _, par := range g.Participants.Map() {
			if !yield(par) {
				return
			}
		}
	}
}

func (g *Game) GetStateParticipantsByState(s participant.State) iter.Seq[*participant.Participant] {
	return func(yield func(*participant.Participant) bool) {
		for _, par := range g.Participants.Map() {
			if par.InState(s) {
				if !yield(par) {
					return
				}
			}
		}
	}
}

func (g *Game) GetParticipantLen() int {
	return g.Participants.Len()
}

func (g *Game) GetParticipantLenByState(s participant.State) int {
	count := 0
	for _, par := range g.Participants.Map() {
		if par.InState(s) {
			count++
		}
	}
	return count
}

func (g *Game) HasEnoughPlayers() bool {
	return g.GetParticipantLen() >= g.Settings.Mode.MinimumTotalPlayers()
}

func (g *Game) IsFull() bool {
	return g.Settings.Mode.MaximumTotalPlayers() != -1 && g.GetParticipantLen() >= g.Settings.Mode.MaximumTotalPlayers()
}

func (g *Game) ParticipantsCallback(fn func(pt *participant.Participant)) {
	for _, pt := range g.Participants.Map() {
		fn(pt)
	}
}

func (g *Game) ParticipantsCallbackByState(s participant.State, fn func(pt *participant.Participant)) {
	for _, par := range g.Participants.Map() {
		if par.InState(s) {
			fn(par)
		}
	}
}

func (g *Game) GetRandomAvailableTeam() (*team.Team, bool) {
	var available []*team.Team
	for _, t := range g.Teams {
		if t.Teammates.Len() < g.Settings.Mode.NumberOfPlayersPerTeam() {
			available = append(available, t)
		}
	}

	if len(available) == 0 {
		return nil, false
	}

	return available[rand.Intn(len(available))], true
}

func (g *Game) GetBalancedAvailableTeam() (*team.Team, bool) {
	var bestTeam *team.Team
	minCount := g.Settings.Mode.NumberOfPlayersPerTeam() + 1

	for _, t := range g.Teams {
		count := t.Teammates.Len()
		if count < g.Settings.Mode.NumberOfPlayersPerTeam() && count < minCount {
			minCount = count
			bestTeam = t
		}
	}

	if bestTeam == nil {
		return nil, false
	}
	return bestTeam, true
}

func (g *Game) AssignTeamToParticipant(pt *participant.Participant, team *team.Team) {
	team.Teammates.Store(pt.Player().UUID(), pt)
}

func (g *Game) RemoveFromTeam(pt *participant.Participant) {
	playerUUID := pt.Player().UUID()
	if t := g.TeamOf(pt); t != nil {
		t.Teammates.Delete(playerUUID)
	}
}

func (g *Game) GetTeamByID(id string) *team.Team {
	for _, t := range g.Teams {
		if t.GetID() == id {
			return t
		}
	}
	return nil
}

func (g *Game) EnemiesOf(pt *participant.Participant) []*participant.Participant {
	var enemies []*participant.Participant

	for other := range g.GetParticipants() {
		if other == pt || g.InSameTeam(pt, other) {
			continue
		}

		enemies = append(enemies, other)
	}

	return enemies
}

func (g *Game) TeamOf(pt *participant.Participant) *team.Team {
	playerUUID := pt.Player().UUID()
	for _, t := range g.Teams {
		if _, ok := t.Teammates.Load(playerUUID); ok {
			return t
		}
	}
	return nil
}

func (g *Game) InSameTeam(a, b *participant.Participant) bool {
	teamA := g.TeamOf(a)
	teamB := g.TeamOf(b)
	return teamA != nil && teamA == teamB
}

func (g *Game) Join(p *player.Player) error {
	if !g.MapLoaded {
		return errors.New("game map is not loaded yet")
	}

	g.Participants.Store(p.UUID(), participant.NewParticipant(p, p.H()))
	return nil
}

func (g *Game) Quit(p *player.Player) {
	g.Participants.Delete(p.UUID())
}
