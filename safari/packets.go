package safari

// Auth / registration packet IDs match ctf-clientside-script (CTFConstants_mem.nut).
const (
	PacketLogin        = 56
	PacketShowLogin    = 57
	PacketHideLogin    = 58
	PacketRegister           = 59
	PacketShowRegister       = 60
	PacketHideRegister       = 61
	PacketRequestRegisterUI  = 94 // client ScriptLoad → server (CTF uses REQUEST_ACC_PREF on load)
	PacketHydraCam           = 95 // server → client hydra camera mode (see hydra_camera.go)
	PacketHydraCamHello      = 96 // client ScriptLoad → server (confirms store/script/main.nut loaded)
)
