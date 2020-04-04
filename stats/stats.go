package stats

import (
	"encoding/json"
)

type Stats struct {
	Hits         int
	Miss         int
	HitCommands  map[string]int
	MissCommands map[string]int
}

var stats *Stats

func init() {
	stats = &Stats{
		Hits:         0,
		Miss:         0,
		HitCommands:  map[string]int{},
		MissCommands: map[string]int{},
	}
}

func Hit() {
	stats.Hits++
}

func Miss() {
	stats.Miss++
}

func Get() *Stats {
	return stats
}

func HitCommand(command string) {
	if _, ok := stats.HitCommands[command]; ok {
		stats.HitCommands[command]++
	} else {
		stats.HitCommands[command] = 1
	}
}

func MissCommand(command string) {
	if _, ok := stats.MissCommands[command]; ok {
		stats.MissCommands[command]++
	} else {
		stats.MissCommands	[command] = 1
	}
}

func (s *Stats) String() string {
	o, e := json.MarshalIndent(s, "", "\t")
	if e != nil {
		panic(e)
	}

	return string(o)
}
