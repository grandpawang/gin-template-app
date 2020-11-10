package redis

import (
	"strings"
)

// redis state
const (
	WatchState = 1 << iota
	MultiState
	SubscribeState
	MonitorState
)

// CommandInfo command info.
type CommandInfo struct {
	Set, Clear int
}

var commandInfos = map[string]CommandInfo{
	"WATCH":      {Set: WatchState},
	"UNWATCH":    {Clear: WatchState},
	"MULTI":      {Set: MultiState},
	"EXEC":       {Clear: WatchState | MultiState},
	"DISCARD":    {Clear: WatchState | MultiState},
	"PSUBSCRIBE": {Set: SubscribeState},
	"SUBSCRIBE":  {Set: SubscribeState},
	"MONITOR":    {Set: MonitorState},
}

func init() {
	for n, ci := range commandInfos {
		commandInfos[strings.ToLower(n)] = ci
	}
}

// LookupCommandInfo get command info.
func LookupCommandInfo(commandName string) CommandInfo {
	if ci, ok := commandInfos[commandName]; ok {
		return ci
	}
	return commandInfos[strings.ToUpper(commandName)]
}
