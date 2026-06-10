package safari

type EventType int

const (
	EvTick EventType = iota
	EvPlayerConnect
	EvPlayerDisconnect
	EvPlayerSpawn
	EvPlayerDeath
	EvPlayerCommand
	EvVehicleExplode
	EvRequestSpawn
	EvClientScriptData
	EvVehicleUpdate
	EvVehicleRespawn
	EvPickupPicked
	EvCheckpointEntered
	EvCheckpointExited
	EvPlayerKeyBind
	EvObjectShot
	EvObjectTouched
	EvPickupRespawn
	EvEntityPoolChange
	EvPlayerUpdate
	EvPlayerEnterVehicle
	EvPlayerExitVehicle
)

type Event struct {
	Type              EventType
	PlayerID          int
	KillerID          int
	Command           string
	VehicleID         int
	PickupID          int
	CheckpointID      int
	KeyBindID         int
	KeyBindReleased   bool
	VehicleUpdateType int
	ObjectID          int
	WeaponID          int
	EntityPool        int
	EntityDeleted     bool
	VehicleSlot       int
	ScriptData        []byte
	SpawnResult       chan bool
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

func NewDeathEvent(playerID, killerID int) Event {
	return Event{Type: EvPlayerDeath, PlayerID: playerID, KillerID: killerID}
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

func NewClientScriptDataEvent(id int, data []byte) Event {
	return Event{Type: EvClientScriptData, PlayerID: id, ScriptData: data}
}

func NewVehicleUpdateEvent(vehicleID, updateType int) Event {
	return Event{Type: EvVehicleUpdate, VehicleID: vehicleID, VehicleUpdateType: updateType}
}

func NewVehicleRespawnEvent(vehicleID int) Event {
	return Event{Type: EvVehicleRespawn, VehicleID: vehicleID}
}

func NewPickupPickedEvent(pickupID, playerID int) Event {
	return Event{Type: EvPickupPicked, PickupID: pickupID, PlayerID: playerID}
}

func NewCheckpointEnteredEvent(checkpointID, playerID int) Event {
	return Event{Type: EvCheckpointEntered, CheckpointID: checkpointID, PlayerID: playerID}
}

func NewCheckpointExitedEvent(checkpointID, playerID int) Event {
	return Event{Type: EvCheckpointExited, CheckpointID: checkpointID, PlayerID: playerID}
}

func NewKeyBindEvent(playerID, bindID int, released bool) Event {
	return Event{Type: EvPlayerKeyBind, PlayerID: playerID, KeyBindID: bindID, KeyBindReleased: released}
}

func NewObjectShotEvent(objectID, playerID, weaponID int) Event {
	return Event{Type: EvObjectShot, ObjectID: objectID, PlayerID: playerID, WeaponID: weaponID}
}

func NewObjectTouchedEvent(objectID, playerID int) Event {
	return Event{Type: EvObjectTouched, ObjectID: objectID, PlayerID: playerID}
}

func NewPickupRespawnEvent(pickupID int) Event {
	return Event{Type: EvPickupRespawn, PickupID: pickupID}
}

func NewEntityPoolChangeEvent(pool, entityID int, deleted bool) Event {
	return Event{Type: EvEntityPoolChange, EntityPool: pool, ObjectID: entityID, EntityDeleted: deleted}
}

func NewPlayerUpdateEvent(playerID, updateType int) Event {
	return Event{Type: EvPlayerUpdate, PlayerID: playerID, VehicleUpdateType: updateType}
}

func NewPlayerEnterVehicleEvent(playerID, vehicleID, slot int) Event {
	return Event{Type: EvPlayerEnterVehicle, PlayerID: playerID, VehicleID: vehicleID, VehicleSlot: slot}
}

func NewPlayerExitVehicleEvent(playerID, vehicleID int) Event {
	return Event{Type: EvPlayerExitVehicle, PlayerID: playerID, VehicleID: vehicleID}
}
