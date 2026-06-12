// Live round HUD (decui).

ScoreboardHUD <- {
	rootId = "safari::hud",
	visible = false,
	res = null,

	function init(res) {
		this.res = res;
	},

	function hide() {
		if (UI.findById(this.rootId) != null) {
			UI.DeleteByID(this.rootId);
		}
		this.visible = false;
	},

	function ensure() {
		if (UI.findById(this.rootId) != null) {
			return;
		}
		local w = this.res.X;

		UI.Canvas({
			id = this.rootId,
			Position = VectorScreen(0, 0),
			Size = VectorScreen(w, 56),
			Colour = Colour(0, 0, 0, 0),
			children = [
				UI.Label({
					id = this.rootId + "::escort",
					Position = VectorScreen(floor(w * 0.06), 10),
					Size = VectorScreen(200, 24),
					Text = "ESCORT 0",
					TextColour = SafariTheme.ESCORT,
					FontSize = 20,
					FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD,
				}),
				UI.Label({
					id = this.rootId + "::timer",
					Position = VectorScreen(floor(w * 0.44), 8),
					Size = VectorScreen(floor(w * 0.12), 28),
					Text = "00:00",
					TextColour = SafariTheme.TEXT,
					FontSize = 22,
					TextAlignment = GUI_ALIGN_CENTERH,
					FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD,
				}),
				UI.Label({
					id = this.rootId + "::defend",
					Position = VectorScreen(floor(w * 0.74), 10),
					Size = VectorScreen(200, 24),
					Text = "DEFEND 0",
					TextColour = SafariTheme.DEFEND,
					FontSize = 20,
					FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD,
				}),
				UI.Label({
					id = this.rootId + "::status",
					Position = VectorScreen(floor(w * 0.2), 40),
					Size = VectorScreen(floor(w * 0.64), 20),
					Text = "",
					TextColour = SafariTheme.MUTED,
					FontSize = 13,
					TextAlignment = GUI_ALIGN_CENTERH,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
			],
		});
		this.visible = true;
	},

	function update(escort, defend, state, mins, secs, hydraHp, cpIdx, cpTotal) {
		if (state == 0) {
			this.hide();
			return;
		}
		this.ensure();

		UI.Label(this.rootId + "::escort").Text = "ESCORT " + escort;
		UI.Label(this.rootId + "::defend").Text = "DEFEND " + defend;

		local mm = mins < 10 ? "0" + mins : mins.tostring();
		local ss = secs < 10 ? "0" + secs : secs.tostring();
		UI.Label(this.rootId + "::timer").Text = mm + ":" + ss;

		local status = "";
		if (state == 2) {
			status = "Round ended";
		} else if (state == 3) {
			status = "Paused";
		} else if (state == 4) {
			status = "Waiting for round start";
		} else if (state == 1) {
			status = "Hydra HP: " + hydraHp.tointeger() + " / 1000";
			if (cpTotal > 0) {
				status += "  |  Checkpoint: " + (cpIdx + 1) + "/" + cpTotal;
			}
		}
		UI.Label(this.rootId + "::status").Text = status;
	},

	function onResize(res) {
		local wasVisible = this.visible;
		this.res = res;
		if (!wasVisible) {
			return;
		}
		this.hide();
	},
};
