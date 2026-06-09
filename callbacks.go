package main

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "plugin.h"
#include <stdlib.h>
#include <string.h>

static void vcmp_copy_player_name(char *dst, size_t buflen, const char *src) {
	if (dst != NULL && buflen > 0) {
		strncpy(dst, src, buflen - 1);
		dst[buflen - 1] = '\0';
	}
}
*/
import "C"

import "unsafe"

//export OnServerInitialise
func OnServerInitialise() C.uint8_t {
	if events.OnServerStart != nil {
		return C.uint8_t(events.OnServerStart())
	}
	return C.uint8_t(FilterAllow)
}

//export OnServerShutdown
func OnServerShutdown() {
	if events.OnServerStop != nil {
		events.OnServerStop()
	}
}

//export OnIncomingConnection
func OnIncomingConnection(playerName *C.char, nameBufferSize C.size_t, userPassword *C.char, ipAddress *C.char) C.uint8_t {
	if events.OnIncomingConnection != nil {
		updated := events.OnIncomingConnection(
			C.GoString(playerName),
			C.GoString(userPassword),
			C.GoString(ipAddress),
		)
		cName := C.CString(updated)
		defer C.free(unsafe.Pointer(cName))
		C.vcmp_copy_player_name(playerName, nameBufferSize, cName)
	}
	return C.uint8_t(FilterAllow)
}

//export OnPlayerConnect
func OnPlayerConnect(playerID C.int32_t) {
	if events.OnPlayerConnect != nil {
		events.OnPlayerConnect(int(playerID))
	}
}

//export OnPlayerDisconnect
func OnPlayerDisconnect(playerID C.int32_t, reason C.vcmpDisconnectReason) {
	if events.OnPlayerDisconnect != nil {
		events.OnPlayerDisconnect(int(playerID), DisconnectReason(reason))
	}
}

//export OnPlayerRequestClass
func OnPlayerRequestClass(playerID C.int32_t, offset C.int32_t) C.uint8_t {
	if events.OnPlayerRequestClass != nil {
		return C.uint8_t(events.OnPlayerRequestClass(int(playerID), int(offset)))
	}
	return C.uint8_t(FilterAllow)
}

//export OnPlayerRequestSpawn
func OnPlayerRequestSpawn(playerID C.int32_t) C.uint8_t {
	if events.OnPlayerRequestSpawn != nil {
		return C.uint8_t(events.OnPlayerRequestSpawn(int(playerID)))
	}
	return C.uint8_t(FilterAllow)
}

//export OnPlayerCommand
func OnPlayerCommand(playerID C.int32_t, message *C.char) C.uint8_t {
	if events.OnPlayerCommand != nil {
		return C.uint8_t(events.OnPlayerCommand(int(playerID), C.GoString(message)))
	}
	return C.uint8_t(FilterAllow)
}

//export OnPlayerEnterVehicle
func OnPlayerEnterVehicle(playerID, vehicleID, slot C.int32_t) {
	if events.OnPlayerEnterVehicle != nil {
		events.OnPlayerEnterVehicle(int(playerID), int(vehicleID), int(slot))
	}
}

//export OnPlayerExitVehicle
func OnPlayerExitVehicle(playerID, vehicleID C.int32_t) {
	if events.OnPlayerExitVehicle != nil {
		events.OnPlayerExitVehicle(int(playerID), int(vehicleID))
	}
}

//export OnVehicleExplode
func OnVehicleExplode(vehicleID C.int32_t) {
	if events.OnVehicleExplode != nil {
		events.OnVehicleExplode(int(vehicleID))
	}
}
