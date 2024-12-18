package noxworld

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
)

func init() {
	ns.Runtime().OnPlayerJoin(playerJoin)
	ns.Runtime().OnPlayerLeave(playerLeave)
	ns.Runtime().OnMapEvent(ns.MapEntry, moveToEntry)
	ns.OnChat(onCommand)
}

func OnFrame() {
}

func onCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil && p == ns.GetHost().Player() {
		data := loadMyNoxWorldData(ns.GetHost().Player())
		switch msg {
		case "-test":
			ns.PrintStrToAll("DEBUG")
			ns.PrintStrToAll(data.Character.LastExitUsed)
		case "-setpos":
			if ns.Waypoint(data.Character.LastExitUsed) != nil {
				ns.GetHost().SetPos(ns.Waypoint(data.Character.LastExitUsed).Pos())
			}
		}
	}
	return msg
}

func checkIfRegistered(p ns.Player) {
	data := loadMyNoxWorldData(p)
	if !data.Character.Registered {
		updateNoxWorldData(ns.GetHost().Player(), func(data *NoxWorldData) {
			data.Character.Registered = true
			ns.AudioEvent(audio.JournalEntryAdd, p.Unit())
			p.PrintStr("Welcome to the open world Nox server. Your character has now been registered!")
			p.PrintStr("Explore the world of Nox your way!")
			p.PrintStr("This server is a work in progress. Please report bugs so the development team can fix them.")
			soulgate := ns.FindClosestObject(p.Unit(), ns.HasTypeName{"SoulGate"})
			if soulgate != nil {
				p.Unit().SetPos(soulgate.Pos())
			}
		})
	} else {
		if ns.Waypoint(data.Character.LastExitUsed) != nil {
			p.Unit().SetPos(ns.Waypoint(data.Character.LastExitUsed).Pos())
			return
		} else {
			soulgate := ns.FindClosestObject(p.Unit(), ns.HasTypeName{"SoulGate"})
			if soulgate != nil {
				p.Unit().SetPos(soulgate.Pos())
			}
		}
	}
}

func moveToEntry() {
	ap := ns.Players()
	for i := 0; i < len(ap); i++ {
		checkIfRegistered(ap[i].Unit().Player())
	}
}

func playerLeave(p ns.Player) {
	// TODO: Store last position.
}

func playerJoin(p ns.Player) bool {
	checkIfRegistered(p)
	return true
}
