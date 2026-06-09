package main

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "plugin.h"
#include <stdlib.h>
#include <string.h>

extern PluginFuncs *g_pf;

static vcmpError vcmp_send_client_script_data(int32_t playerId, const void *data, size_t size) {
	if (g_pf && g_pf->SendClientScriptData) return g_pf->SendClientScriptData(playerId, data, size);
	return vcmpErrorNullArgument;
}
static vcmpError vcmp_send_game_message(int32_t playerId, int32_t type, const char *msg) {
	if (g_pf && g_pf->SendGameMessage) return g_pf->SendGameMessage(playerId, type, "%s", msg);
	return vcmpErrorNullArgument;
}
static vcmpError vcmp_add_radio_stream(int32_t radioId, const char *name, const char *url, uint8_t listed) {
	if (g_pf && g_pf->AddRadioStream) return g_pf->AddRadioStream(radioId, name, url, listed);
	return vcmpErrorNullArgument;
}
static vcmpError vcmp_remove_radio_stream(int32_t radioId) {
	if (g_pf && g_pf->RemoveRadioStream) return g_pf->RemoveRadioStream(radioId);
	return vcmpErrorNoSuchEntity;
}
static uint8_t vcmp_is_player_streamed(int32_t checked, int32_t forPlayer) {
	if (g_pf && g_pf->IsPlayerStreamedForPlayer) return g_pf->IsPlayerStreamedForPlayer(checked, forPlayer);
	return 0;
}
static uint8_t vcmp_is_vehicle_streamed(int32_t vehicleId, int32_t playerId) {
	if (g_pf && g_pf->IsVehicleStreamedForPlayer) return g_pf->IsVehicleStreamedForPlayer(vehicleId, playerId);
	return 0;
}
static uint8_t vcmp_is_pickup_streamed(int32_t pickupId, int32_t playerId) {
	if (g_pf && g_pf->IsPickupStreamedForPlayer) return g_pf->IsPickupStreamedForPlayer(pickupId, playerId);
	return 0;
}
static uint8_t vcmp_is_object_streamed(int32_t objectId, int32_t playerId) {
	if (g_pf && g_pf->IsObjectStreamedForPlayer) return g_pf->IsObjectStreamedForPlayer(objectId, playerId);
	return 0;
}
static uint8_t vcmp_is_checkpoint_streamed(int32_t cpId, int32_t playerId) {
	if (g_pf && g_pf->IsCheckPointStreamedForPlayer) return g_pf->IsCheckPointStreamedForPlayer(cpId, playerId);
	return 0;
}
static vcmpError vcmp_set_player_weapon(int32_t playerId, int32_t weaponId, int32_t ammo) {
	if (g_pf && g_pf->SetPlayerWeapon) return g_pf->SetPlayerWeapon(playerId, weaponId, ammo);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_get_player_weapon(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerWeapon) return g_pf->GetPlayerWeapon(playerId);
	return 0;
}
static int32_t vcmp_get_player_weapon_ammo(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerWeaponAmmo) return g_pf->GetPlayerWeaponAmmo(playerId);
	return 0;
}
static vcmpError vcmp_set_player_weapon_slot(int32_t playerId, int32_t slot) {
	if (g_pf && g_pf->SetPlayerWeaponSlot) return g_pf->SetPlayerWeaponSlot(playerId, slot);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_get_player_weapon_slot(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerWeaponSlot) return g_pf->GetPlayerWeaponSlot(playerId);
	return 0;
}
static int32_t vcmp_get_player_weapon_at_slot(int32_t playerId, int32_t slot) {
	if (g_pf && g_pf->GetPlayerWeaponAtSlot) return g_pf->GetPlayerWeaponAtSlot(playerId, slot);
	return 0;
}
static int32_t vcmp_get_player_ammo_at_slot(int32_t playerId, int32_t slot) {
	if (g_pf && g_pf->GetPlayerAmmoAtSlot) return g_pf->GetPlayerAmmoAtSlot(playerId, slot);
	return 0;
}
static vcmpError vcmp_remove_player_weapon(int32_t playerId, int32_t weaponId) {
	if (g_pf && g_pf->RemovePlayerWeapon) return g_pf->RemovePlayerWeapon(playerId, weaponId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_weapon_data_value(int32_t weaponId, int32_t fieldId, double value) {
	if (g_pf && g_pf->SetWeaponDataValue) return g_pf->SetWeaponDataValue(weaponId, fieldId, value);
	return vcmpErrorArgumentOutOfBounds;
}
static double vcmp_get_weapon_data_value(int32_t weaponId, int32_t fieldId) {
	if (g_pf && g_pf->GetWeaponDataValue) return g_pf->GetWeaponDataValue(weaponId, fieldId);
	return 0;
}
static vcmpError vcmp_set_player_health(int32_t playerId, float health) {
	if (g_pf && g_pf->SetPlayerHealth) return g_pf->SetPlayerHealth(playerId, health);
	return vcmpErrorNoSuchEntity;
}
static float vcmp_get_player_health(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerHealth) return g_pf->GetPlayerHealth(playerId);
	return 0;
}
static vcmpError vcmp_set_player_armour(int32_t playerId, float armour) {
	if (g_pf && g_pf->SetPlayerArmour) return g_pf->SetPlayerArmour(playerId, armour);
	return vcmpErrorNoSuchEntity;
}
static float vcmp_get_player_armour(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerArmour) return g_pf->GetPlayerArmour(playerId);
	return 0;
}
static vcmpError vcmp_give_player_money(int32_t playerId, int32_t amount) {
	if (g_pf && g_pf->GivePlayerMoney) return g_pf->GivePlayerMoney(playerId, amount);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_player_money(int32_t playerId, int32_t amount) {
	if (g_pf && g_pf->SetPlayerMoney) return g_pf->SetPlayerMoney(playerId, amount);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_get_player_money(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerMoney) return g_pf->GetPlayerMoney(playerId);
	return 0;
}
static vcmpError vcmp_set_player_position(int32_t playerId, float x, float y, float z) {
	if (g_pf && g_pf->SetPlayerPosition) return g_pf->SetPlayerPosition(playerId, x, y, z);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_player_heading(int32_t playerId, float angle) {
	if (g_pf && g_pf->SetPlayerHeading) return g_pf->SetPlayerHeading(playerId, angle);
	return vcmpErrorNoSuchEntity;
}
static float vcmp_get_player_heading(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerHeading) return g_pf->GetPlayerHeading(playerId);
	return 0;
}
static vcmpError vcmp_set_player_name(int32_t playerId, const char *name) {
	if (g_pf && g_pf->SetPlayerName) return g_pf->SetPlayerName(playerId, name);
	return vcmpErrorNoSuchEntity;
}
static vcmpPlayerState vcmp_get_player_state(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerState) return g_pf->GetPlayerState(playerId);
	return vcmpPlayerStateNone;
}
static vcmpError vcmp_set_player_option(int32_t playerId, vcmpPlayerOption option, uint8_t toggle) {
	if (g_pf && g_pf->SetPlayerOption) return g_pf->SetPlayerOption(playerId, option, toggle);
	return vcmpErrorNoSuchEntity;
}
static uint8_t vcmp_get_player_option(int32_t playerId, vcmpPlayerOption option) {
	if (g_pf && g_pf->GetPlayerOption) return g_pf->GetPlayerOption(playerId, option);
	return 0;
}
static uint8_t vcmp_is_player_spawned(int32_t playerId) {
	if (g_pf && g_pf->IsPlayerSpawned) return g_pf->IsPlayerSpawned(playerId);
	return 0;
}
static vcmpError vcmp_force_player_spawn(int32_t playerId) {
	if (g_pf && g_pf->ForcePlayerSpawn) return g_pf->ForcePlayerSpawn(playerId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_force_player_select(int32_t playerId) {
	if (g_pf && g_pf->ForcePlayerSelect) return g_pf->ForcePlayerSelect(playerId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_kick_player(int32_t playerId) {
	if (g_pf && g_pf->KickPlayer) return g_pf->KickPlayer(playerId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_ban_player(int32_t playerId) {
	if (g_pf && g_pf->BanPlayer) return g_pf->BanPlayer(playerId);
	return vcmpErrorNoSuchEntity;
}
static void vcmp_get_player_ip(int32_t playerId, char *buf, size_t buflen) {
	if (buf && buflen > 0) buf[0] = '\0';
	if (g_pf && g_pf->GetPlayerIP && buf && buflen > 0) g_pf->GetPlayerIP(playerId, buf, buflen);
}
static int32_t vcmp_get_player_ping(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerPing) return g_pf->GetPlayerPing(playerId);
	return 0;
}
static vcmpError vcmp_remove_player_from_vehicle(int32_t playerId) {
	if (g_pf && g_pf->RemovePlayerFromVehicle) return g_pf->RemovePlayerFromVehicle(playerId);
	return vcmpErrorNoSuchEntity;
}
static vcmpPlayerVehicle vcmp_get_player_in_vehicle_status(int32_t playerId) {
	if (g_pf && g_pf->GetPlayerInVehicleStatus) return g_pf->GetPlayerInVehicleStatus(playerId);
	return vcmpPlayerVehicleOut;
}
static vcmpError vcmp_respawn_vehicle(int32_t vehicleId) {
	if (g_pf && g_pf->RespawnVehicle) return g_pf->RespawnVehicle(vehicleId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_explode_vehicle(int32_t vehicleId) {
	if (g_pf && g_pf->ExplodeVehicle) return g_pf->ExplodeVehicle(vehicleId);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_vehicle_option(int32_t vehicleId, vcmpVehicleOption option, uint8_t toggle) {
	if (g_pf && g_pf->SetVehicleOption) return g_pf->SetVehicleOption(vehicleId, option, toggle);
	return vcmpErrorNoSuchEntity;
}
static uint8_t vcmp_get_vehicle_option(int32_t vehicleId, vcmpVehicleOption option) {
	if (g_pf && g_pf->GetVehicleOption) return g_pf->GetVehicleOption(vehicleId, option);
	return 0;
}
static vcmpError vcmp_set_vehicle_world(int32_t vehicleId, int32_t world) {
	if (g_pf && g_pf->SetVehicleWorld) return g_pf->SetVehicleWorld(vehicleId, world);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_get_vehicle_world(int32_t vehicleId) {
	if (g_pf && g_pf->GetVehicleWorld) return g_pf->GetVehicleWorld(vehicleId);
	return 0;
}
static int32_t vcmp_get_vehicle_model(int32_t vehicleId) {
	if (g_pf && g_pf->GetVehicleModel) return g_pf->GetVehicleModel(vehicleId);
	return 0;
}
static vcmpError vcmp_set_vehicle_rotation_euler(int32_t vehicleId, float x, float y, float z) {
	if (g_pf && g_pf->SetVehicleRotationEuler) return g_pf->SetVehicleRotationEuler(vehicleId, x, y, z);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_vehicle_colour(int32_t vehicleId, int32_t c1, int32_t c2) {
	if (g_pf && g_pf->SetVehicleColour) return g_pf->SetVehicleColour(vehicleId, c1, c2);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_set_vehicle_radio(int32_t vehicleId, int32_t radioId) {
	if (g_pf && g_pf->SetVehicleRadio) return g_pf->SetVehicleRadio(vehicleId, radioId);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_get_vehicle_radio(int32_t vehicleId) {
	if (g_pf && g_pf->GetVehicleRadio) return g_pf->GetVehicleRadio(vehicleId);
	return 0;
}
static uint8_t vcmp_is_vehicle_wrecked(int32_t vehicleId) {
	if (g_pf && g_pf->IsVehicleWrecked) return g_pf->IsVehicleWrecked(vehicleId);
	return 0;
}
static int32_t vcmp_create_pickup(int32_t model, int32_t world, int32_t qty, float x, float y, float z, int32_t alpha, uint8_t automatic) {
	if (g_pf && g_pf->CreatePickup) return g_pf->CreatePickup(model, world, qty, x, y, z, alpha, automatic);
	return -1;
}
static vcmpError vcmp_delete_pickup(int32_t pickupId) {
	if (g_pf && g_pf->DeletePickup) return g_pf->DeletePickup(pickupId);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_create_object(int32_t model, int32_t world, float x, float y, float z, int32_t alpha) {
	if (g_pf && g_pf->CreateObject) return g_pf->CreateObject(model, world, x, y, z, alpha);
	return -1;
}
static vcmpError vcmp_delete_object(int32_t objectId) {
	if (g_pf && g_pf->DeleteObject) return g_pf->DeleteObject(objectId);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_create_checkpoint(int32_t playerId, int32_t world, uint8_t sphere, float x, float y, float z, int32_t r, int32_t g, int32_t b, int32_t a, float radius) {
	if (g_pf && g_pf->CreateCheckPoint) return g_pf->CreateCheckPoint(playerId, world, sphere, x, y, z, r, g, b, a, radius);
	return -1;
}
static vcmpError vcmp_delete_checkpoint(int32_t cpId) {
	if (g_pf && g_pf->DeleteCheckPoint) return g_pf->DeleteCheckPoint(cpId);
	return vcmpErrorNoSuchEntity;
}
static int32_t vcmp_create_coord_blip(int32_t index, int32_t world, float x, float y, float z, int32_t scale, uint32_t colour, int32_t sprite) {
	if (g_pf && g_pf->CreateCoordBlip) return g_pf->CreateCoordBlip(index, world, x, y, z, scale, colour, sprite);
	return -1;
}
static vcmpError vcmp_destroy_coord_blip(int32_t index) {
	if (g_pf && g_pf->DestroyCoordBlip) return g_pf->DestroyCoordBlip(index);
	return vcmpErrorNoSuchEntity;
}
static vcmpError vcmp_create_explosion(int32_t world, int32_t type, float x, float y, float z, int32_t responsible, uint8_t atGround) {
	if (g_pf && g_pf->CreateExplosion) return g_pf->CreateExplosion(world, type, x, y, z, responsible, atGround);
	return vcmpErrorArgumentOutOfBounds;
}
static void vcmp_set_hour(int32_t hour) {
	if (g_pf && g_pf->SetHour) g_pf->SetHour(hour);
}
static int32_t vcmp_get_hour(void) {
	if (g_pf && g_pf->GetHour) return g_pf->GetHour();
	return 0;
}
static void vcmp_set_weather(int32_t weather) {
	if (g_pf && g_pf->SetWeather) g_pf->SetWeather(weather);
}
static int32_t vcmp_get_weather(void) {
	if (g_pf && g_pf->GetWeather) return g_pf->GetWeather();
	return 0;
}
static void vcmp_set_gravity(float gravity) {
	if (g_pf && g_pf->SetGravity) g_pf->SetGravity(gravity);
}
static float vcmp_get_gravity(void) {
	if (g_pf && g_pf->GetGravity) return g_pf->GetGravity();
	return 0;
}
static vcmpError vcmp_register_key_bind(int32_t bindId, uint8_t onRelease, int32_t k1, int32_t k2, int32_t k3) {
	if (g_pf && g_pf->RegisterKeyBind) return g_pf->RegisterKeyBind(bindId, onRelease, k1, k2, k3);
	return vcmpErrorArgumentOutOfBounds;
}
static vcmpError vcmp_remove_key_bind(int32_t bindId) {
	if (g_pf && g_pf->RemoveKeyBind) return g_pf->RemoveKeyBind(bindId);
	return vcmpErrorNoSuchEntity;
}
static uint8_t vcmp_check_entity_exists(vcmpEntityPool pool, int32_t index) {
	if (g_pf && g_pf->CheckEntityExists) return g_pf->CheckEntityExists(pool, index);
	return 0;
}
*/
import "C"

import "unsafe"

func boolToU8(v bool) C.uint8_t {
	if v {
		return 1
	}
	return 0
}

func bridgeSendClientScriptData(playerID int, data []byte) error {
	if len(data) == 0 {
		return bridgeError(C.vcmp_send_client_script_data(C.int32_t(playerID), nil, 0))
	}
	return bridgeError(C.vcmp_send_client_script_data(C.int32_t(playerID), unsafe.Pointer(&data[0]), C.size_t(len(data))))
}

func bridgeSendGameMessage(playerID, msgType int, msg string) error {
	c := cString(msg)
	defer freeCString(c)
	return bridgeError(C.vcmp_send_game_message(C.int32_t(playerID), C.int32_t(msgType), c))
}

func bridgeAddRadioStream(radioID int, name, url string, listed bool) error {
	cName, cURL := cString(name), cString(url)
	defer freeCString(cName)
	defer freeCString(cURL)
	return bridgeError(C.vcmp_add_radio_stream(C.int32_t(radioID), cName, cURL, boolToU8(listed)))
}

func bridgeRemoveRadioStream(radioID int) error {
	return bridgeError(C.vcmp_remove_radio_stream(C.int32_t(radioID)))
}

func bridgeIsPlayerStreamedForPlayer(checkedPlayerID, forPlayerID int) bool {
	return C.vcmp_is_player_streamed(C.int32_t(checkedPlayerID), C.int32_t(forPlayerID)) != 0
}

func bridgeIsVehicleStreamedForPlayer(vehicleID, playerID int) bool {
	return C.vcmp_is_vehicle_streamed(C.int32_t(vehicleID), C.int32_t(playerID)) != 0
}

func bridgeIsPickupStreamedForPlayer(pickupID, playerID int) bool {
	return C.vcmp_is_pickup_streamed(C.int32_t(pickupID), C.int32_t(playerID)) != 0
}

func bridgeIsObjectStreamedForPlayer(objectID, playerID int) bool {
	return C.vcmp_is_object_streamed(C.int32_t(objectID), C.int32_t(playerID)) != 0
}

func bridgeIsCheckPointStreamedForPlayer(checkpointID, playerID int) bool {
	return C.vcmp_is_checkpoint_streamed(C.int32_t(checkpointID), C.int32_t(playerID)) != 0
}

func bridgeSetPlayerWeapon(playerID, weaponID, ammo int) error {
	return bridgeError(C.vcmp_set_player_weapon(C.int32_t(playerID), C.int32_t(weaponID), C.int32_t(ammo)))
}

func bridgeGetPlayerWeapon(playerID int) int {
	return int(C.vcmp_get_player_weapon(C.int32_t(playerID)))
}

func bridgeGetPlayerWeaponAmmo(playerID int) int {
	return int(C.vcmp_get_player_weapon_ammo(C.int32_t(playerID)))
}

func bridgeSetPlayerWeaponSlot(playerID, slot int) error {
	return bridgeError(C.vcmp_set_player_weapon_slot(C.int32_t(playerID), C.int32_t(slot)))
}

func bridgeGetPlayerWeaponSlot(playerID int) int {
	return int(C.vcmp_get_player_weapon_slot(C.int32_t(playerID)))
}

func bridgeGetPlayerWeaponAtSlot(playerID, slot int) int {
	return int(C.vcmp_get_player_weapon_at_slot(C.int32_t(playerID), C.int32_t(slot)))
}

func bridgeGetPlayerAmmoAtSlot(playerID, slot int) int {
	return int(C.vcmp_get_player_ammo_at_slot(C.int32_t(playerID), C.int32_t(slot)))
}

func bridgeRemovePlayerWeapon(playerID, weaponID int) error {
	return bridgeError(C.vcmp_remove_player_weapon(C.int32_t(playerID), C.int32_t(weaponID)))
}

func bridgeSetWeaponDataValue(weaponID, fieldID int, value float64) error {
	return bridgeError(C.vcmp_set_weapon_data_value(C.int32_t(weaponID), C.int32_t(fieldID), C.double(value)))
}

func bridgeGetWeaponDataValue(weaponID, fieldID int) float64 {
	return float64(C.vcmp_get_weapon_data_value(C.int32_t(weaponID), C.int32_t(fieldID)))
}

func bridgeSetPlayerHealth(playerID int, health float32) error {
	return bridgeError(C.vcmp_set_player_health(C.int32_t(playerID), C.float(health)))
}

func bridgeGetPlayerHealth(playerID int) float32 {
	return float32(C.vcmp_get_player_health(C.int32_t(playerID)))
}

func bridgeSetPlayerArmour(playerID int, armour float32) error {
	return bridgeError(C.vcmp_set_player_armour(C.int32_t(playerID), C.float(armour)))
}

func bridgeGetPlayerArmour(playerID int) float32 {
	return float32(C.vcmp_get_player_armour(C.int32_t(playerID)))
}

func bridgeGivePlayerMoney(playerID, amount int) error {
	return bridgeError(C.vcmp_give_player_money(C.int32_t(playerID), C.int32_t(amount)))
}

func bridgeSetPlayerMoney(playerID, amount int) error {
	return bridgeError(C.vcmp_set_player_money(C.int32_t(playerID), C.int32_t(amount)))
}

func bridgeGetPlayerMoney(playerID int) int {
	return int(C.vcmp_get_player_money(C.int32_t(playerID)))
}

func bridgeSetPlayerPosition(playerID int, pos Vec3) error {
	return bridgeError(C.vcmp_set_player_position(C.int32_t(playerID), C.float(pos.X), C.float(pos.Y), C.float(pos.Z)))
}

func bridgeSetPlayerHeading(playerID int, angle float32) error {
	return bridgeError(C.vcmp_set_player_heading(C.int32_t(playerID), C.float(angle)))
}

func bridgeGetPlayerHeading(playerID int) float32 {
	return float32(C.vcmp_get_player_heading(C.int32_t(playerID)))
}

func bridgeSetPlayerName(playerID int, name string) error {
	c := cString(name)
	defer freeCString(c)
	return bridgeError(C.vcmp_set_player_name(C.int32_t(playerID), c))
}

func bridgeGetPlayerState(playerID int) PlayerState {
	return PlayerState(C.vcmp_get_player_state(C.int32_t(playerID)))
}

func bridgeSetPlayerOption(playerID int, option PlayerOption, on bool) error {
	return bridgeError(C.vcmp_set_player_option(C.int32_t(playerID), C.vcmpPlayerOption(option), boolToU8(on)))
}

func bridgeGetPlayerOption(playerID int, option PlayerOption) bool {
	return C.vcmp_get_player_option(C.int32_t(playerID), C.vcmpPlayerOption(option)) != 0
}

func bridgeIsPlayerSpawned(playerID int) bool {
	return C.vcmp_is_player_spawned(C.int32_t(playerID)) != 0
}

func bridgeForcePlayerSpawn(playerID int) error {
	return bridgeError(C.vcmp_force_player_spawn(C.int32_t(playerID)))
}

func bridgeForcePlayerSelect(playerID int) error {
	return bridgeError(C.vcmp_force_player_select(C.int32_t(playerID)))
}

func bridgeKickPlayer(playerID int) error {
	return bridgeError(C.vcmp_kick_player(C.int32_t(playerID)))
}

func bridgeBanPlayer(playerID int) error {
	return bridgeError(C.vcmp_ban_player(C.int32_t(playerID)))
}

func bridgeGetPlayerIP(playerID int) string {
	buf := (*[64]C.char)(C.malloc(64))
	defer C.free(unsafe.Pointer(buf))
	C.vcmp_get_player_ip(C.int32_t(playerID), &buf[0], 64)
	return C.GoString(&buf[0])
}

func bridgeGetPlayerPing(playerID int) int {
	return int(C.vcmp_get_player_ping(C.int32_t(playerID)))
}

func bridgeRemovePlayerFromVehicle(playerID int) error {
	return bridgeError(C.vcmp_remove_player_from_vehicle(C.int32_t(playerID)))
}

func bridgeGetPlayerInVehicleStatus(playerID int) PlayerVehicle {
	return PlayerVehicle(C.vcmp_get_player_in_vehicle_status(C.int32_t(playerID)))
}

func bridgeRespawnVehicle(vehicleID int) error {
	return bridgeError(C.vcmp_respawn_vehicle(C.int32_t(vehicleID)))
}

func bridgeExplodeVehicle(vehicleID int) error {
	return bridgeError(C.vcmp_explode_vehicle(C.int32_t(vehicleID)))
}

func bridgeSetVehicleOption(vehicleID int, option VehicleOption, on bool) error {
	return bridgeError(C.vcmp_set_vehicle_option(C.int32_t(vehicleID), C.vcmpVehicleOption(option), boolToU8(on)))
}

func bridgeGetVehicleOption(vehicleID int, option VehicleOption) bool {
	return C.vcmp_get_vehicle_option(C.int32_t(vehicleID), C.vcmpVehicleOption(option)) != 0
}

func bridgeSetVehicleWorld(vehicleID, world int) error {
	return bridgeError(C.vcmp_set_vehicle_world(C.int32_t(vehicleID), C.int32_t(world)))
}

func bridgeGetVehicleWorld(vehicleID int) int {
	return int(C.vcmp_get_vehicle_world(C.int32_t(vehicleID)))
}

func bridgeGetVehicleModel(vehicleID int) int {
	return int(C.vcmp_get_vehicle_model(C.int32_t(vehicleID)))
}

func bridgeSetVehicleRotationEuler(vehicleID int, rot Vec3) error {
	return bridgeError(C.vcmp_set_vehicle_rotation_euler(C.int32_t(vehicleID), C.float(rot.X), C.float(rot.Y), C.float(rot.Z)))
}

func bridgeSetVehicleColour(vehicleID, primary, secondary int) error {
	return bridgeError(C.vcmp_set_vehicle_colour(C.int32_t(vehicleID), C.int32_t(primary), C.int32_t(secondary)))
}

func bridgeSetVehicleRadio(vehicleID, radioID int) error {
	return bridgeError(C.vcmp_set_vehicle_radio(C.int32_t(vehicleID), C.int32_t(radioID)))
}

func bridgeGetVehicleRadio(vehicleID int) int {
	return int(C.vcmp_get_vehicle_radio(C.int32_t(vehicleID)))
}

func bridgeIsVehicleWrecked(vehicleID int) bool {
	return C.vcmp_is_vehicle_wrecked(C.int32_t(vehicleID)) != 0
}

func bridgeCreatePickup(model, world, quantity int, pos Vec3, alpha int, automatic bool) int {
	return int(C.vcmp_create_pickup(
		C.int32_t(model), C.int32_t(world), C.int32_t(quantity),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z),
		C.int32_t(alpha), boolToU8(automatic),
	))
}

func bridgeDeletePickup(pickupID int) error {
	return bridgeError(C.vcmp_delete_pickup(C.int32_t(pickupID)))
}

func bridgeCreateObject(model, world int, pos Vec3, alpha int) int {
	return int(C.vcmp_create_object(C.int32_t(model), C.int32_t(world), C.float(pos.X), C.float(pos.Y), C.float(pos.Z), C.int32_t(alpha)))
}

func bridgeDeleteObject(objectID int) error {
	return bridgeError(C.vcmp_delete_object(C.int32_t(objectID)))
}

func bridgeCreateCheckPoint(playerID, world int, sphere bool, pos Vec3, r, g, b, alpha int, radius float32) int {
	return int(C.vcmp_create_checkpoint(
		C.int32_t(playerID), C.int32_t(world), boolToU8(sphere),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z),
		C.int32_t(r), C.int32_t(g), C.int32_t(b), C.int32_t(alpha), C.float(radius),
	))
}

func bridgeDeleteCheckPoint(checkpointID int) error {
	return bridgeError(C.vcmp_delete_checkpoint(C.int32_t(checkpointID)))
}

func bridgeCreateCoordBlip(index, world int, pos Vec3, scale int, colour uint32, sprite int) int {
	return int(C.vcmp_create_coord_blip(
		C.int32_t(index), C.int32_t(world),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z),
		C.int32_t(scale), C.uint32_t(colour), C.int32_t(sprite),
	))
}

func bridgeDestroyCoordBlip(index int) error {
	return bridgeError(C.vcmp_destroy_coord_blip(C.int32_t(index)))
}

func bridgeCreateExplosion(world, explosionType int, pos Vec3, responsiblePlayerID int, atGroundLevel bool) error {
	return bridgeError(C.vcmp_create_explosion(
		C.int32_t(world), C.int32_t(explosionType),
		C.float(pos.X), C.float(pos.Y), C.float(pos.Z),
		C.int32_t(responsiblePlayerID), boolToU8(atGroundLevel),
	))
}

func bridgeSetHour(hour int) { C.vcmp_set_hour(C.int32_t(hour)) }
func bridgeGetHour() int      { return int(C.vcmp_get_hour()) }
func bridgeSetWeather(w int)  { C.vcmp_set_weather(C.int32_t(w)) }
func bridgeGetWeather() int   { return int(C.vcmp_get_weather()) }
func bridgeSetGravity(g float32) { C.vcmp_set_gravity(C.float(g)) }
func bridgeGetGravity() float32  { return float32(C.vcmp_get_gravity()) }

func bridgeRegisterKeyBind(bindID int, onRelease bool, keyOne, keyTwo, keyThree int) error {
	code := C.vcmp_register_key_bind(C.int32_t(bindID), boolToU8(onRelease), C.int32_t(keyOne), C.int32_t(keyTwo), C.int32_t(keyThree))
	return bridgeError(C.vcmpError(code))
}

func bridgeRemoveKeyBind(bindID int) error {
	return bridgeError(C.vcmp_remove_key_bind(C.int32_t(bindID)))
}

func bridgeCheckEntityExists(pool EntityPool, index int) bool {
	return C.vcmp_check_entity_exists(C.vcmpEntityPool(pool), C.int32_t(index)) != 0
}

func bridgeError(code C.vcmpError) error {
	if code == C.vcmpErrorNone {
		return nil
	}
	return vcmpError(code)
}

type vcmpError C.vcmpError

func (e vcmpError) Error() string {
	switch C.vcmpError(e) {
	case C.vcmpErrorNone:
		return "none"
	case C.vcmpErrorNoSuchEntity:
		return "no such entity"
	case C.vcmpErrorBufferTooSmall:
		return "buffer too small"
	case C.vcmpErrorTooLargeInput:
		return "too large input"
	case C.vcmpErrorArgumentOutOfBounds:
		return "argument out of bounds"
	case C.vcmpErrorNullArgument:
		return "null argument"
	case C.vcmpErrorPoolExhausted:
		return "pool exhausted"
	case C.vcmpErrorInvalidName:
		return "invalid name"
	case C.vcmpErrorRequestDenied:
		return "request denied"
	default:
		return "unknown vcmp error"
	}
}
