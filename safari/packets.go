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
	PacketScoreboard         = 97 // server → client round HUD
	PacketHydraCamCycle      = 98 // client → server request camera cycle (H key)
	PacketSelectPack           = 99 // client → server pack selection
	PacketShowPacks            = 100 // server → client open pack picker
	PacketHidePacks            = 101 // server → client close pack picker
	PacketRoundEndStats        = 102 // server → client end-of-round scoreboard
	PacketRequestShowPacks     = 103 // client → server open pack picker (P key)
	PacketLobbyLeaderboard     = 104 // server → client 3D lobby leaderboard boards
)
