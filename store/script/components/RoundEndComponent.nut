// Round-end fullscreen overlay (decui).

class RoundEndComponent extends Component {
	className = "RoundEndComponent";

	winnerTeam = 0;
	escortScore = 0;
	defendScore = 0;
	reason = "";
	res = null;
	rootId = null;

	constructor(o) {
		this.id = o.id;
		this.rootId = o.id + "::overlay";
		this.winnerTeam = o.winnerTeam;
		this.escortScore = o.escortScore;
		this.defendScore = o.defendScore;
		this.reason = o.reason;
		this.res = o.res;
		base.constructor(this.id, o);
		this.metadata.list = "canvas";
		this.build();
	}

	function build() {
		local winnerName = winnerTeam == Teams.ESCORT ? "ESCORT" : "DEFENDERS";
		local winnerColour = winnerTeam == Teams.ESCORT ? SafariTheme.ESCORT : SafariTheme.DEFEND;

		local panel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(520, 120),
			Colour = SafariTheme.PANEL,
			children = [
				UI.Label({
					id = id + "::banner",
					align = "center",
					Position = VectorScreen(0, 20),
					Size = VectorScreen(520, 32),
					Text = winnerName + " WIN - " + escortScore + " vs " + defendScore,
					TextColour = winnerColour,
					FontSize = 24,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::reason",
					align = "center",
					Position = VectorScreen(0, 54),
					Size = VectorScreen(520, 22),
					Text = reason,
					TextColour = SafariTheme.TEXT,
					FontSize = 14,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::hint",
					Position = VectorScreen(380, 88),
					Size = VectorScreen(120, 18),
					Text = "[P] Close",
					TextColour = SafariTheme.MUTED,
					FontSize = 12,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
			],
		});
		panel.addBorders({ size = 2, color = winnerColour });

		UI.Canvas({
			id = rootId,
			align = "center",
			Size = VectorScreen(res.X, res.Y),
			Position = VectorScreen(0, 0),
			Colour = SafariTheme.OVERLAY,
			children = [panel],
		});

		GUI.SetMouseEnabled(true);
	}

	function destroy() {
		if (UI.findById(rootId) != null) {
			UI.DeleteByID(rootId);
		}
	}
}
