package main

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "plugin.h"
*/
import "C"

const (
	MaxPlayers = 100
	PluginName = "Safari"

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

type Quat struct {
	X float32
	Y float32
	Z float32
	W float32
}

type ServerOption int

const (
	ServerOptionSyncFrameLimiter       ServerOption = ServerOption(C.vcmpServerOptionSyncFrameLimiter)
	ServerOptionFrameLimiter           ServerOption = ServerOption(C.vcmpServerOptionFrameLimiter)
	ServerOptionTaxiBoostJump          ServerOption = ServerOption(C.vcmpServerOptionTaxiBoostJump)
	ServerOptionDriveOnWater           ServerOption = ServerOption(C.vcmpServerOptionDriveOnWater)
	ServerOptionFastSwitch             ServerOption = ServerOption(C.vcmpServerOptionFastSwitch)
	ServerOptionFriendlyFire           ServerOption = ServerOption(C.vcmpServerOptionFriendlyFire)
	ServerOptionDisableDriveBy         ServerOption = ServerOption(C.vcmpServerOptionDisableDriveBy)
	ServerOptionPerfectHandling        ServerOption = ServerOption(C.vcmpServerOptionPerfectHandling)
	ServerOptionFlyingCars             ServerOption = ServerOption(C.vcmpServerOptionFlyingCars)
	ServerOptionJumpSwitch             ServerOption = ServerOption(C.vcmpServerOptionJumpSwitch)
	ServerOptionShowMarkers            ServerOption = ServerOption(C.vcmpServerOptionShowMarkers)
	ServerOptionOnlyShowTeamMarkers    ServerOption = ServerOption(C.vcmpServerOptionOnlyShowTeamMarkers)
	ServerOptionStuntBike              ServerOption = ServerOption(C.vcmpServerOptionStuntBike)
	ServerOptionShootInAir             ServerOption = ServerOption(C.vcmpServerOptionShootInAir)
	ServerOptionShowNameTags           ServerOption = ServerOption(C.vcmpServerOptionShowNameTags)
	ServerOptionJoinMessages           ServerOption = ServerOption(C.vcmpServerOptionJoinMessages)
	ServerOptionDeathMessages          ServerOption = ServerOption(C.vcmpServerOptionDeathMessages)
	ServerOptionChatTagsEnabled        ServerOption = ServerOption(C.vcmpServerOptionChatTagsEnabled)
	ServerOptionUseClasses             ServerOption = ServerOption(C.vcmpServerOptionUseClasses)
	ServerOptionWallGlitch             ServerOption = ServerOption(C.vcmpServerOptionWallGlitch)
	ServerOptionDisableBackfaceCulling ServerOption = ServerOption(C.vcmpServerOptionDisableBackfaceCulling)
	ServerOptionDisableHeliBladeDamage ServerOption = ServerOption(C.vcmpServerOptionDisableHeliBladeDamage)
)

type PlayerOption int

const (
	PlayerOptionControllable     PlayerOption = PlayerOption(C.vcmpPlayerOptionControllable)
	PlayerOptionDriveBy          PlayerOption = PlayerOption(C.vcmpPlayerOptionDriveBy)
	PlayerOptionWhiteScanlines   PlayerOption = PlayerOption(C.vcmpPlayerOptionWhiteScanlines)
	PlayerOptionGreenScanlines   PlayerOption = PlayerOption(C.vcmpPlayerOptionGreenScanlines)
	PlayerOptionWidescreen       PlayerOption = PlayerOption(C.vcmpPlayerOptionWidescreen)
	PlayerOptionShowMarkers      PlayerOption = PlayerOption(C.vcmpPlayerOptionShowMarkers)
	PlayerOptionCanAttack        PlayerOption = PlayerOption(C.vcmpPlayerOptionCanAttack)
	PlayerOptionHasMarker        PlayerOption = PlayerOption(C.vcmpPlayerOptionHasMarker)
	PlayerOptionChatTagsEnabled  PlayerOption = PlayerOption(C.vcmpPlayerOptionChatTagsEnabled)
	PlayerOptionDrunkEffects     PlayerOption = PlayerOption(C.vcmpPlayerOptionDrunkEffects)
)

type VehicleOption int

const (
	VehicleOptionDoorsLocked VehicleOption = VehicleOption(C.vcmpVehicleOptionDoorsLocked)
	VehicleOptionAlarm       VehicleOption = VehicleOption(C.vcmpVehicleOptionAlarm)
	VehicleOptionLights      VehicleOption = VehicleOption(C.vcmpVehicleOptionLights)
	VehicleOptionRadioLocked VehicleOption = VehicleOption(C.vcmpVehicleOptionRadioLocked)
	VehicleOptionGhost       VehicleOption = VehicleOption(C.vcmpVehicleOptionGhost)
	VehicleOptionSiren       VehicleOption = VehicleOption(C.vcmpVehicleOptionSiren)
	VehicleOptionSingleUse   VehicleOption = VehicleOption(C.vcmpVehicleOptionSingleUse)
)

type PickupOption int

const (
	PickupOptionSingleUse PickupOption = PickupOption(C.vcmpPickupOptionSingleUse)
)

type DisconnectReason int

const (
	DisconnectTimeout   DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonTimeout)
	DisconnectQuit      DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonQuit)
	DisconnectKick      DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonKick)
	DisconnectCrash     DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonCrash)
	DisconnectAntiCheat DisconnectReason = DisconnectReason(C.vcmpDisconnectReasonAntiCheat)
)

type BodyPart int

const (
	BodyPartBody      BodyPart = BodyPart(C.vcmpBodyPartBody)
	BodyPartTorso     BodyPart = BodyPart(C.vcmpBodyPartTorso)
	BodyPartLeftArm   BodyPart = BodyPart(C.vcmpBodyPartLeftArm)
	BodyPartRightArm  BodyPart = BodyPart(C.vcmpBodyPartRightArm)
	BodyPartLeftLeg   BodyPart = BodyPart(C.vcmpBodyPartLeftLeg)
	BodyPartRightLeg  BodyPart = BodyPart(C.vcmpBodyPartRightLeg)
	BodyPartHead      BodyPart = BodyPart(C.vcmpBodyPartHead)
	BodyPartInVehicle BodyPart = BodyPart(C.vcmpBodyPartInVehicle)
)

type PlayerState int

const (
	PlayerStateNone           PlayerState = PlayerState(C.vcmpPlayerStateNone)
	PlayerStateNormal         PlayerState = PlayerState(C.vcmpPlayerStateNormal)
	PlayerStateAim            PlayerState = PlayerState(C.vcmpPlayerStateAim)
	PlayerStateDriver         PlayerState = PlayerState(C.vcmpPlayerStateDriver)
	PlayerStatePassenger      PlayerState = PlayerState(C.vcmpPlayerStatePassenger)
	PlayerStateEnterDriver    PlayerState = PlayerState(C.vcmpPlayerStateEnterDriver)
	PlayerStateEnterPassenger PlayerState = PlayerState(C.vcmpPlayerStateEnterPassenger)
	PlayerStateExit           PlayerState = PlayerState(C.vcmpPlayerStateExit)
	PlayerStateUnspawned      PlayerState = PlayerState(C.vcmpPlayerStateUnspawned)
)

type PlayerUpdate int

const (
	PlayerUpdateNormal    PlayerUpdate = PlayerUpdate(C.vcmpPlayerUpdateNormal)
	PlayerUpdateAiming    PlayerUpdate = PlayerUpdate(C.vcmpPlayerUpdateAiming)
	PlayerUpdateDriver    PlayerUpdate = PlayerUpdate(C.vcmpPlayerUpdateDriver)
	PlayerUpdatePassenger PlayerUpdate = PlayerUpdate(C.vcmpPlayerUpdatePassenger)
)

type PlayerVehicle int

const (
	PlayerVehicleOut      PlayerVehicle = PlayerVehicle(C.vcmpPlayerVehicleOut)
	PlayerVehicleEntering PlayerVehicle = PlayerVehicle(C.vcmpPlayerVehicleEntering)
	PlayerVehicleExiting  PlayerVehicle = PlayerVehicle(C.vcmpPlayerVehicleExiting)
	PlayerVehicleIn       PlayerVehicle = PlayerVehicle(C.vcmpPlayerVehicleIn)
)

type VehicleUpdate int

const (
	VehicleUpdateDriverSync VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdateDriverSync)
	VehicleUpdateOtherSync  VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdateOtherSync)
	VehicleUpdatePosition   VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdatePosition)
	VehicleUpdateHealth     VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdateHealth)
	VehicleUpdateColour     VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdateColour)
	VehicleUpdateRotation   VehicleUpdate = VehicleUpdate(C.vcmpVehicleUpdateRotation)
)

type EntityPool int

const (
	EntityPoolVehicle    EntityPool = EntityPool(C.vcmpEntityPoolVehicle)
	EntityPoolObject     EntityPool = EntityPool(C.vcmpEntityPoolObject)
	EntityPoolPickup     EntityPool = EntityPool(C.vcmpEntityPoolPickup)
	EntityPoolRadio      EntityPool = EntityPool(C.vcmpEntityPoolRadio)
	EntityPoolBlip       EntityPool = EntityPool(C.vcmpEntityPoolBlip)
	EntityPoolCheckPoint EntityPool = EntityPool(C.vcmpEntityPoolCheckPoint)
)

type FilterResult uint8

const (
	FilterAllow FilterResult = 1
	FilterDeny  FilterResult = 0
)
