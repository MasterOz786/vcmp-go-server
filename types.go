package main

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "plugin.h"
*/
import "C"

const (
	MaxPlayers = 100
	PluginName = "GoDemo"

	ColourWhite     uint32 = 0xFFFFFFFF
	ColourYellowish uint32 = 0xFFFF5500
	ColourGreen     uint32 = 0xFF00FF00
	ColourYellow    uint32 = 0xFFFFFF00
	ColourRed       uint32 = 0xFFFF4040
	ColourCyan      uint32 = 0xFF00FFFF
	ColourResponse  uint32 = 0xFFC8FF96
	ColourTimer     uint32 = 0xFFFFC8C8
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type ServerOption int

const (
	ServerOptionJoinMessages  ServerOption = ServerOption(C.vcmpServerOptionJoinMessages)
	ServerOptionDeathMessages ServerOption = ServerOption(C.vcmpServerOptionDeathMessages)
	ServerOptionUseClasses    ServerOption = ServerOption(C.vcmpServerOptionUseClasses)
)

type DisconnectReason int

const (
	DisconnectTimeout   DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonTimeout)
	DisconnectQuit      DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonQuit)
	DisconnectKick      DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonKick)
	DisconnectCrash     DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonCrash)
	DisconnectAntiCheat DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonAntiCheat)
)

type FilterResult uint8

const (
	FilterAllow FilterResult = 1
	FilterDeny  FilterResult = 0
)
