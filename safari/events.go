package safari

type EventType int

const (
	EvTick EventType = iota
	EvPlayerConnect
	EvPlayerDisconnect
	EvPlayerSpawn
	EvPlayerCommand
	EvVehicleExplode
	EvRequestSpawn
)

type Event struct {
	Type     EventType
	PlayerID int
	Command  string
	VehicleID int
	SpawnResult chan bool
}

func NewTickEvent() Event {
	return Event{Type: EvTick}
}

func NewConnectEvent(id int) Event {
	return Event{Type: EvPlayerConnect, PlayerID: id}
}

func NewDisconnectEvent(id int) Event {
	return Event{Type: EvPlayerDisconnect, PlayerID: id}
}

func NewSpawnEvent(id int) Event {
	return Event{Type: EvPlayerSpawn, PlayerID: id}
}

func NewCommandEvent(id int, cmd string) Event {
	return Event{Type: EvPlayerCommand, PlayerID: id, Command: cmd}
}

func NewVehicleExplodeEvent(vid int) Event {
	return Event{Type: EvVehicleExplode, VehicleID: vid}
}

func NewRequestSpawnEvent(id int) Event {
	return Event{Type: EvRequestSpawn, PlayerID: id, SpawnResult: make(chan bool, 1)}
}
