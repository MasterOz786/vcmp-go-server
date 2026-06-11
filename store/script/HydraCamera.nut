// Hydra camera: server applies SetCamera on H (/hydraview). Client only requests cycles.

SafariHydraCam <- {
	mode = HydraCam.OFF,
	vehicleId = -1,
	hydraKey = null,

	names = ["Pilot (default)", "Chase cam", "Side orbit", "Tactical overhead"],

	function init() {
		this.hydraKey = KeyBind(0x48); // H — also registered server-side as bind 4
		print("[safari] client ready — H cycles Hydra camera via server");
	},

	function sendHello() {
		local s = Stream();
		s.WriteInt(Packets.HYDRA_CAM_HELLO);
		Server.SendData(s);
	},

	function requestCycle() {
		local s = Stream();
		s.WriteInt(Packets.HYDRA_CAM_CYCLE);
		Server.SendData(s);
	},

	function applyMode(mode, vid) {
		if (vid >= 0) {
			this.vehicleId = vid;
		}
		this.mode = mode;
		if (mode >= HydraCam.DEFAULT && mode < this.names.len()) {
			print("[safari] hydra view: " + this.names[mode]);
		}
	}
};
