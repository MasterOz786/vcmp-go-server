// Client-side Hydra follow camera. Server sends vehicle id + mode; math runs here.
// Use H (not V) — V is native GTA vehicle camera and fights custom views.

SafariHydraCam <- {
	mode = HydraCam.OFF,
	vehicleId = -1,
	smoothX = 0.0,
	smoothY = 0.0,
	smoothZ = 0.0,
	hasSmooth = false,
	lastHeading = 90.0,
	lerp = 0.18,
	lastTick = 0,
	intervalMs = 33,
	hydraKey = null,

	names = ["Pilot (default)", "Chase cam", "Side orbit", "Tactical overhead"],

	function init() {
		this.hydraKey = KeyBind(0x48); // H
		print("[safari] hydra camera client loaded — press H or /hydraview (F8 console)");
	},

	function sendHello() {
		local s = Stream();
		s.WriteInt(Packets.HYDRA_CAM_HELLO);
		Server.SendData(s);
	},

	function applyMode(mode, vid) {
		if (vid >= 0) {
			this.vehicleId = vid;
		}
		if (mode == HydraCam.OFF) {
			this.mode = HydraCam.OFF;
			this.hasSmooth = false;
			this.restorePilotCamera();
			return;
		}
		this.mode = mode;
		this.hasSmooth = false;
		if (mode <= HydraCam.DEFAULT) {
			this.restorePilotCamera();
		}
	},

	function restorePilotCamera() {
		local plr = World.FindLocalPlayer();
		if (plr) {
			plr.RestoreCamera();
		}
	},

	function cycleMode() {
		local veh = this.resolveVehicle();
		if (!veh) {
			print("[safari] hydra cam: enter a Hydra first (/testhydra)");
			return;
		}
		this.mode = (this.mode + 1) % 4;
		if (this.mode < 0) {
			this.mode = 0;
		}
		this.hasSmooth = false;
		if (this.mode <= HydraCam.DEFAULT) {
			this.restorePilotCamera();
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
		case HydraCam.CHASE:
			cam.X = pos.X - sinH * 38.0;
			cam.Y = pos.Y + cosH * 38.0;
			cam.Z = pos.Z + 14.0;
			break;
		case HydraCam.SIDE:
			cam.X = pos.X + cosH * 28.0;
			cam.Y = pos.Y + sinH * 28.0;
			cam.Z = pos.Z + 10.0;
			break;
		case HydraCam.TACTICAL:
			cam.X = pos.X;
			cam.Y = pos.Y;
			cam.Z = pos.Z + 55.0;
			break;
		default:
			return null;
		}
		return { pos = cam, look = look };
	},

	// Trust vehicle id from server; do not require model 6460 (fallback heli still gets a cam).
	function resolveVehicle() {
		if (this.vehicleId < 0) {
			return null;
		}
		local veh = World.FindVehicle(this.vehicleId);
		if (!veh || veh.Health <= 0.0) {
			return null;
		}
		local plr = World.FindLocalPlayer();
		if (!plr) {
			return null;
		}
		local driver = veh.GetOccupant(0);
		if (driver != plr) {
			return null;
		}
		return veh;
	},

	function onScriptProcess() {
		local veh = this.resolveVehicle();
		if (!veh) {
			if (this.mode > HydraCam.DEFAULT) {
				this.applyMode(HydraCam.OFF, -1);
			}
			return;
		}

		if (this.mode <= HydraCam.DEFAULT) {
			return;
		}

		local now = System.GetTimestamp();
		if (this.lastTick != 0 && (now - this.lastTick) < this.intervalMs) {
			return;
		}
		this.lastTick = now;

		local plr = World.FindLocalPlayer();
		if (!plr) {
			return;
		}

		local off = this.cameraOffsets(veh, this.mode);
		if (off == null) {
			return;
		}

		if (!this.hasSmooth) {
			this.smoothX = off.pos.X;
			this.smoothY = off.pos.Y;
			this.smoothZ = off.pos.Z;
			this.hasSmooth = true;
		} else {
			this.smoothX += (off.pos.X - this.smoothX) * this.lerp;
			this.smoothY += (off.pos.Y - this.smoothY) * this.lerp;
			this.smoothZ += (off.pos.Z - this.smoothZ) * this.lerp;
		}

		plr.SetCameraPos(
			Vector(this.smoothX, this.smoothY, this.smoothZ),
			off.look
		);
	}
};
