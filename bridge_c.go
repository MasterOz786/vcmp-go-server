package main

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "plugin.h"
#include <stdlib.h>
#include <string.h>

static PluginFuncs *g_pf;
static void vcmp_set_funcs(PluginFuncs *pf) { g_pf = pf; }

static void vcmp_log_msg(const char *msg) {
	if (g_pf && g_pf->LogMessage) g_pf->LogMessage("%s", msg);
}
static uint64_t vcmp_get_time(void) {
	return g_pf && g_pf->GetTime ? g_pf->GetTime() : 0;
}
static void vcmp_set_server_name(const char *name) {
	if (g_pf && g_pf->SetServerName) g_pf->SetServerName(name);
}
static void vcmp_get_server_name(char *buf, size_t buflen) {
	if (buf && buflen > 0) buf[0] = '\0';
	if (g_pf && g_pf->GetServerName && buf && buflen > 0) g_pf->GetServerName(buf, buflen);
}
static void vcmp_set_gamemode_text(const char *text) {
	if (g_pf && g_pf->SetGameModeText) g_pf->SetGameModeText(text);
}
static void vcmp_set_server_option(vcmpServerOption option, uint8_t toggle) {
	if (g_pf && g_pf->SetServerOption) g_pf->SetServerOption(option, toggle);
}
static void vcmp_set_spawn_pos(float x, float y, float z) {
	if (g_pf && g_pf->SetSpawnPlayerPosition) g_pf->SetSpawnPlayerPosition(x, y, z);
}
static int32_t vcmp_add_player_class(int32_t teamId, uint32_t colour, int32_t modelIndex, float x, float y, float z, float angle, int32_t w1, int32_t w1a, int32_t w2, int32_t w2a, int32_t w3, int32_t w3a) {
	if (g_pf && g_pf->AddPlayerClass) return g_pf->AddPlayerClass(teamId, colour, modelIndex, x, y, z, angle, w1, w1a, w2, w2a, w3, w3a);
	return -1;
}
static void vcmp_send_client_message(int32_t playerId, uint32_t colour, const char *msg) {
	if (g_pf && g_pf->SendClientMessage) g_pf->SendClientMessage(playerId, colour, "%s", msg);
}
static void vcmp_get_player_name(int32_t playerId, char *buf, size_t buflen) {
	if (buf && buflen > 0) buf[0] = '\0';
	if (g_pf && g_pf->GetPlayerName && buf && buflen > 0) g_pf->GetPlayerName(playerId, buf, (int32_t)buflen);
}
static int32_t vcmp_get_player_id_from_name(const char *name) {
	if (g_pf && g_pf->GetPlayerIdFromName) return g_pf->GetPlayerIdFromName(name);
	return -1;
}
static uint8_t vcmp_is_player_connected(int32_t playerId) {
	if (g_pf && g_pf->IsPlayerConnected) return g_pf->IsPlayerConnected(playerId);
	return 0;
}
static uint8_t vcmp_is_player_admin(int32_t playerId) {
	if (g_pf && g_pf->IsPlayerAdmin) return g_pf->IsPlayerAdmin(playerId);
	return 0;
}
static void vcmp_set_player_admin(int32_t playerId, uint8_t toggle) {
	if (g_pf && g_pf->SetPlayerAdmin) g_pf->SetPlayerAdmin(playerId, toggle);
}
static void vcmp_get_player_position(int32_t playerId, float *x, float *y, float *z) {
	if (g_pf && g_pf->GetPlayerPosition && x && y && z) g_pf->GetPlayerPosition(playerId, x, y, z);
}
static int32_t vcmp_get_player_world(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerWorld) return g_pf->GetPlayerWorld(playerId);
	return 0;
}
static void vcmp_set_player_world(int32_t playerId, int32_t world) {
	if (g_pf && g_pf->SetPlayerWorld) g_pf->SetPlayerWorld(playerId, world);
}
static int32_t vcmp_get_player_vehicle_id(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerVehicleId) return g_pf->GetPlayerVehicleId(playerId);
	return -1;
}
static void vcmp_put_player_in_vehicle(int32_t playerId, int32_t vehicleId, int32_t slot, uint8_t makeRoom, uint8_t warp) {
	if (g_pf && g_pf->PutPlayerInVehicle) g_pf->PutPlayerInVehicle(playerId, vehicleId, slot, makeRoom, warp);
}
static int32_t vcmp_create_vehicle(int32_t model, int32_t world, float x, float y, float z, float angle, int32_t c1, int32_t c2) {
	if (g_pf && g_pf->CreateVehicle) return g_pf->CreateVehicle(model, world, x, y, z, angle, c1, c2);
	return -1;
}
static void vcmp_get_vehicle_position(int32_t vehicleId, float *x, float *y, float *z) {
	if (g_pf && g_pf->GetVehiclePosition && x && y && z) g_pf->GetVehiclePosition(vehicleId, x, y, z);
}
static float vcmp_get_vehicle_health(int32_t vehicleId) {
	if (g_pf && g_pf->GetVehicleHealth) return g_pf->GetVehicleHealth(vehicleId);
	return 0;
}
static int32_t vcmp_get_vehicle_occupant(int32_t vehicleId, int32_t slot) {
	if (g_pf && g_pf->GetVehicleOccupant) return g_pf->GetVehicleOccupant(vehicleId, slot);
	return -1;
}
static void vcmp_set_vehicle_part_status(int32_t vehicleId, int32_t partId, int32_t status) {
	if (g_pf && g_pf->SetVehiclePartStatus) g_pf->SetVehiclePartStatus(vehicleId, partId, status);
}
static void vcmp_set_vehicle_tyre_status(int32_t vehicleId, int32_t tyreId, int32_t status) {
	if (g_pf && g_pf->SetVehicleTyreStatus) g_pf->SetVehicleTyreStatus(vehicleId, tyreId, status);
}
*/
import "C"

import "unsafe"

func bindPluginAPI(pf *C.PluginFuncs) { C.vcmp_set_funcs(pf) }

func cString(s string) *C.char { return C.CString(s) }
func freeCString(p *C.char)    { C.free(unsafe.Pointer(p)) }

func bridgeLog(msg string) {
	c := cString(msg)
	defer freeCString(c)
	C.vcmp_log_msg(c)
}

func bridgeSetServerName(name string) {
	c := cString(name)
	defer freeCString(c)
	C.vcmp_set_server_name(c)
}

func bridgeGetServerName() string {
	buf := (*[128]C.char)(C.malloc(128))
	defer C.free(unsafe.Pointer(buf))
	C.vcmp_get_server_name(&buf[0], 128)
	return C.GoString(&buf[0])
}

func bridgeSetGameModeText(text string) {
	c := cString(text)
	defer freeCString(c)
	C.vcmp_set_gamemode_text(c)
}

func bridgeSetServerOption(option ServerOption, on bool) {
	t := C.uint8_t(0)
	if on {
		t = 1
	}
	C.vcmp_set_server_option(C.vcmpServerOption(option), t)
}

func bridgeSetSpawnPos(pos Vec3) {
	C.vcmp_set_spawn_pos(C.float(pos.X), C.float(pos.Y), C.float(pos.Z))
}

func bridgeAddPlayerClass(teamID int, colour uint32, model int, pos Vec3, angle float32, w []int) int {
	return int(C.vcmp_add_player_class(
		C.int32_t(teamID), C.uint32_t(colour), C.int32_t(model),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z), C.float(angle),
		C.int32_t(w[0]), C.int32_t(w[1]), C.int32_t(w[2]), C.int32_t(w[3]), C.int32_t(w[4]), C.int32_t(w[5]),
	))
}

func bridgeBroadcast(colour uint32, msg string) {
	for id := 0; id < MaxPlayers; id++ {
		if bridgeIsConnected(id) {
			bridgeSendClientMessage(id, colour, msg)
		}
	}
}

func bridgeSendClientMessage(playerID int, colour uint32, msg string) {
	c := cString(msg)
	defer freeCString(c)
	C.vcmp_send_client_message(C.int32_t(playerID), C.uint32_t(colour), c)
}

func bridgePlayerName(playerID int) string {
	buf := (*[128]C.char)(C.malloc(128))
	defer C.free(unsafe.Pointer(buf))
	C.vcmp_get_player_name(C.int32_t(playerID), &buf[0], 128)
	return C.GoString(&buf[0])
}

func bridgePlayerIDFromName(name string) int {
	c := cString(name)
	defer freeCString(c)
	return int(C.vcmp_get_player_id_from_name(c))
}

func bridgeIsConnected(playerID int) bool {
	return C.vcmp_is_player_connected(C.int32_t(playerID)) != 0
}

func bridgeIsAdmin(playerID int) bool {
	return C.vcmp_is_player_admin(C.int32_t(playerID)) != 0
}

func bridgeSetAdmin(playerID int, admin bool) {
	t := C.uint8_t(0)
	if admin {
		t = 1
	}
	C.vcmp_set_player_admin(C.int32_t(playerID), t)
}

func bridgePlayerPos(playerID int) Vec3 {
	var x, y, z C.float
	C.vcmp_get_player_position(C.int32_t(playerID), &x, &y, &z)
	return Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
}

func bridgePlayerWorld(playerID int) int {
	return int(C.vcmp_get_player_world(C.int32_t(playerID)))
}

func bridgeSetPlayerWorld(playerID, world int) {
	C.vcmp_set_player_world(C.int32_t(playerID), C.int32_t(world))
}

func bridgePlayerVehicleID(playerID int) int {
	return int(C.vcmp_get_player_vehicle_id(C.int32_t(playerID)))
}

func bridgePutInVehicle(playerID, vehicleID, slot int, makeRoom, warp bool) {
	mr, w := C.uint8_t(0), C.uint8_t(0)
	if makeRoom {
		mr = 1
	}
	if warp {
		w = 1
	}
	C.vcmp_put_player_in_vehicle(C.int32_t(playerID), C.int32_t(vehicleID), C.int32_t(slot), mr, w)
}

func bridgeCreateVehicle(model, world int, pos Vec3, angle float32, c1, c2 int) int {
	return int(C.vcmp_create_vehicle(
		C.int32_t(model), C.int32_t(world),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z), C.float(angle),
		C.int32_t(c1), C.int32_t(c2),
	))
}

func bridgeVehiclePos(vehicleID int) Vec3 {
	var x, y, z C.float
	C.vcmp_get_vehicle_position(C.int32_t(vehicleID), &x, &y, &z)
	return Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
}

func bridgeVehicleHealth(vehicleID int) float32 {
	return float32(C.vcmp_get_vehicle_health(C.int32_t(vehicleID)))
}

func bridgeVehicleOccupant(vehicleID, slot int) int {
	return int(C.vcmp_get_vehicle_occupant(C.int32_t(vehicleID), C.int32_t(slot)))
}

func bridgeBreakVehicle(vehicleID int) {
	C.vcmp_set_vehicle_part_status(C.int32_t(vehicleID), 0, 3)
	C.vcmp_set_vehicle_tyre_status(C.int32_t(vehicleID), 3, 1)
	C.vcmp_set_vehicle_part_status(C.int32_t(vehicleID), 4, 2)
}
