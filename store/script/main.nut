// Project Safari — client-side Hydra camera (store/script/main.nut).
// Server sends vehicle id via SendClientScriptData; camera math runs here (not on server).
// Use H (not V) — V is the native GTA vehicle camera toggle and causes conflict/lag.

const PACKET_HYDRA_CAM = 95;
const PACKET_HYDRA_HELLO = 96;

SafariHydraCam <- {
	HYDRA_MODEL = 6460,
	mode = -1,
	vehicleId = -1,
	smoothX = 0.0,
	smoothY = 0.0,
	smoothZ = 0.0,
	hasSmooth = false,
	lastHeading = 90.0,
	lerp = 0.15,
	lastTick = 0,
	intervalMs = 33,
	hydraKey = null,

	names = ["Pilot (default)", "Chase cam", "Side orbit", "Tactical overhead"],

	function sendToServer(pkt) {
		local s = Stream();
		s.WriteInt(pkt);
		SendServerData(s);
	},

	function onServerData(stream) {
		local mode = stream.ReadInt();
		local vid = stream.ReadInt();
		this.applyMode(mode, vid);
	},

	function applyMode(mode, vid) {
		if (vid >= 0) {
			this.vehicleId = vid;
		}
		if (mode == -1) {
			this.mode = -1;
			this.hasSmooth = false;
			local plr = World.FindLocalPlayer();
			if (plr) plr.RestoreCamera();
			return;
		}
		this.mode = mode;
		this.hasSmooth = false;
		local plr = World.FindLocalPlayer();
		if (plr && mode <= 0) {
			plr.RestoreCamera();
		}
	},

	function cycleMode() {
		local veh = this.resolveVehicle();
		if (!veh) {
			print("[safari] hydra cam: not in hydra");
			return;
		}
		this.mode = (this.mode + 1) % 4;
		if (this.mode < 0) this.mode = 0;
		this.hasSmooth = false;
		if (this.mode <= 0) {
			local plr = World.FindLocalPlayer();
			if (plr) plr.RestoreCamera();
		}
		print("[safari] hydra view: " + this.names[this.mode]);
	},

	function vehicleHeading(veh) {
		local spd = Vector(veh.Speed);
		if ((spd.X * spd.X + spd.Y * spd.Y) > 0.25) {
			this.lastHeading = atan2(spd.Y, spd.X) * 57.2957795;
		}
		return this.lastHeading;
	},

	function cameraOffsets(veh, mode) {
		local pos = Vector(veh.Position);
		local look = Vector(pos.X, pos.Y, pos.Z + 1.0);
		local heading = this.vehicleHeading(veh);
		local rad = (heading + 90.0) * 0.0174532925;
		local sinH = sin(rad);
		local cosH = cos(rad);
		local cam = Vector(pos.X, pos.Y, pos.Z);

		switch (mode) {
		case 1:
			cam.X = pos.X - sinH * 38.0;
			cam.Y = pos.Y + cosH * 38.0;
			cam.Z = pos.Z + 14.0;
			break;
		case 2:
			cam.X = pos.X + cosH * 28.0;
			cam.Y = pos.Y + sinH * 28.0;
			cam.Z = pos.Z + 10.0;
			break;
		case 3:
			cam.X = pos.X;
			cam.Y = pos.Y;
			cam.Z = pos.Z + 55.0;
			break;
		default:
			return null;
		}
		return { pos = cam, look = look };
	},

	function resolveVehicle() {
		if (this.vehicleId >= 0) {
			local veh = World.FindVehicle(this.vehicleId);
			if (veh && veh.ModelIndex == this.HYDRA_MODEL) {
				return veh;
			}
		}
		return null;
	},

	function onScriptProcess() {
		local veh = this.resolveVehicle();
		if (!veh) {
			if (this.mode >= 0) {
				this.applyMode(-1, -1);
			}
			return;
		}

		if (this.mode <= 0) {
			return;
		}

		local now = System.GetTimestamp();
		if (this.lastTick != 0 && (now - this.lastTick) < this.intervalMs) {
			return;
		}
		this.lastTick = now;

		local plr = World.FindLocalPlayer();
		if (!plr) return;

		local off = this.cameraOffsets(veh, this.mode);
		if (off == null) return;

		local tx = off.pos.X;
		local ty = off.pos.Y;
		local tz = off.pos.Z;

		if (!this.hasSmooth) {
			this.smoothX = tx;
			this.smoothY = ty;
			this.smoothZ = tz;
			this.hasSmooth = true;
		} else {
			this.smoothX += (tx - this.smoothX) * this.lerp;
			this.smoothY += (ty - this.smoothY) * this.lerp;
			this.smoothZ += (tz - this.smoothZ) * this.lerp;
		}

		plr.SetCameraPos(
			Vector(this.smoothX, this.smoothY, this.smoothZ),
			off.look
		);
	}
};

function SafariHydraCam_HandleServerStream(stream) {
	local pkt = stream.ReadInt();
	if (pkt == PACKET_HYDRA_CAM) {
		SafariHydraCam.onServerData(stream);
	}
}

function Script::ScriptLoad() {
	SafariHydraCam.hydraKey = KeyBind(0x48); // H — must init here, not at table creation
	print("[safari] hydra camera client loaded — use H or /hydraview (this message is in F8, not the server console)");
	try {
		SafariHydraCam.sendToServer(PACKET_HYDRA_HELLO);
	} catch (e) {
		print("[safari] hydra cam hello to server failed: " + e);
	}
}

function Script::ScriptProcess() {
	SafariHydraCam.onScriptProcess();
}

function Server::ServerData(stream) {
	SafariHydraCam_HandleServerStream(stream);
}

function KeyBind::OnDown(bind) {
	if (SafariHydraCam.hydraKey != null && bind == SafariHydraCam.hydraKey) {
		SafariHydraCam.cycleMode();
	}
}
