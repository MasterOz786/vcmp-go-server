package stream

// Auth / registration packet IDs match ctf-clientside-script (CTFConstants_mem.nut).
const (
	PacketLogin              = 56
	PacketShowLogin          = 57
	PacketHideLogin          = 58
	PacketRegister           = 59
	PacketShowRegister       = 60
	PacketHideRegister       = 61
	PacketRequestRegisterUI  = 94 // client ScriptLoad → server
	PacketHydraCam           = 95 // server → client hydra camera mode
	PacketHydraCamHello      = 96 // client ScriptLoad → server
	PacketScoreboard         = 97 // server → client round HUD
	PacketHydraCamCycle      = 98 // client → server request camera cycle
	PacketSelectPack         = 99 // client → server pack selection
	PacketShowPacks          = 100 // server → client open pack picker
	PacketHidePacks          = 101 // server → client close pack picker
	PacketRoundEndStats      = 102 // server → client end-of-round scoreboard
	PacketRequestShowPacks   = 103 // client → server open pack picker
	PacketLobbyLeaderboard   = 104 // server → client 3D lobby leaderboard boards
	PacketPackFeedback       = 105 // server → client pack UI status
	PacketRequestHideLeaderboard = 106 // client → server close leaderboard UI
)
