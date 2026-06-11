class StreamsController {
	windows = null;

	constructor(windows) {
		this.windows = windows;
	}

	function process(stream) {
		local type = stream.ReadInt();

		if (type == Packets.HIDE_REGISTER) {
			if (windows.registerWindow != null && windows.registerWindow.registerBtn != null) {
				GUI.SetFocusedElement(windows.registerWindow.registerBtn);
				windows.registerWindow.clear();
			}
		} else if (type == Packets.SHOW_REGISTER) {
			windows.registerWindow.createWindow();
		} else if (type == Packets.HYDRA_CAM) {
			SafariHydraCam.applyMode(stream.ReadInt(), stream.ReadInt());
		} else if (type == Packets.SCOREBOARD) {
			scoreboard.update(
				stream.ReadInt(),
				stream.ReadInt(),
				stream.ReadInt(),
				stream.ReadInt(),
				stream.ReadInt(),
				stream.ReadFloat(),
				stream.ReadInt(),
				stream.ReadInt()
			);
		} else if (type == Packets.SHOW_PACKS) {
			local team = stream.ReadInt();
			local currentPack = stream.ReadInt();
			windows.packsWindow.createWindow(team, currentPack);
		} else if (type == Packets.HIDE_PACKS) {
			windows.packsWindow.clear();
		} else if (type == Packets.ROUND_END_STATS) {
			local winnerTeam = stream.ReadInt();
			if (winnerTeam < 0) {
				windows.roundScoreboard.hide();
				return;
			}
			local escortScore = stream.ReadInt();
			local defendScore = stream.ReadInt();
			local reason = stream.ReadString();
			local count = stream.ReadInt();
			local players = [];
			for (local i = 0; i < count; i++) {
				players.append({
					name = stream.ReadString(),
					team = stream.ReadInt(),
					points = stream.ReadInt(),
					kills = stream.ReadInt(),
					deaths = stream.ReadInt(),
				});
			}
			windows.roundScoreboard.show(winnerTeam, escortScore, defendScore, reason, players);
		}
	}
}
