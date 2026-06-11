// Round scoreboard HUD — server pushes PacketScoreboard every second during rounds.

ScoreboardHUD <- {
	escortLabel = null,
	defendLabel = null,
	timerLabel = null,
	statusLabel = null,
	visible = false,
	res = null,

	function init(res) {
		this.res = res;
	},

	function hide() {
		if (this.escortLabel != null) {
			this.escortLabel.Detach();
			this.escortLabel = null;
		}
		if (this.defendLabel != null) {
			this.defendLabel.Detach();
			this.defendLabel = null;
		}
		if (this.timerLabel != null) {
			this.timerLabel.Detach();
			this.timerLabel = null;
		}
		if (this.statusLabel != null) {
			this.statusLabel.Detach();
			this.statusLabel = null;
		}
		this.visible = false;
	},

	function ensureLabels() {
		if (this.escortLabel != null) {
			return;
		}
		local w = this.res.X;
		local escortColour = Colour(100, 220, 255);
		local defendColour = Colour(255, 115, 115);

		this.escortLabel = GUILabel(VectorScreen(floor(w * 0.06), 10), escortColour, "ESCORT 0");
		this.escortLabel.FontSize = 20;
		this.escortLabel.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;

		this.defendLabel = GUILabel(VectorScreen(floor(w * 0.74), 10), defendColour, "DEFEND 0");
		this.defendLabel.FontSize = 20;
		this.defendLabel.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;

		this.timerLabel = GUILabel(VectorScreen(floor(w * 0.44), 8), Colour(245, 248, 255), "00:00");
		this.timerLabel.FontSize = 22;
		this.timerLabel.TextAlignment = GUI_ALIGN_CENTERH;
		this.timerLabel.Size = VectorScreen(floor(w * 0.12), 28);
		this.timerLabel.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;

		this.statusLabel = GUILabel(VectorScreen(floor(w * 0.2), 40), Colour(170, 180, 200), "");
		this.statusLabel.FontSize = 13;
		this.statusLabel.Size = VectorScreen(floor(w * 0.64), 24);
		this.statusLabel.TextAlignment = GUI_ALIGN_CENTERH;
		this.statusLabel.FontFlags = GUI_FFLAG_OUTLINE;
		this.visible = true;
	},

	function update(escort, defend, state, mins, secs, hydraHp, cpIdx, cpTotal) {
		if (state == 0) {
			this.hide();
			return;
		}
		this.ensureLabels();

		this.escortLabel.Text = "ESCORT " + escort;
		this.defendLabel.Text = "DEFEND " + defend;

		local mm = mins < 10 ? "0" + mins : mins.tostring();
		local ss = secs < 10 ? "0" + secs : secs.tostring();
		this.timerLabel.Text = mm + ":" + ss;

		if (state == 2) {
			this.statusLabel.Text = "Round ended";
		} else if (state == 3) {
			this.statusLabel.Text = "Paused";
		} else if (state == 4) {
			this.statusLabel.Text = "Waiting for round start";
		} else if (state == 1) {
			local status = "Hydra HP: " + hydraHp.tointeger() + " / 1000";
			if (cpTotal > 0) {
				status += "  |  Checkpoint: " + (cpIdx + 1) + "/" + cpTotal;
			}
			this.statusLabel.Text = status;
		}
	},

	function onResize(res) {
		local wasVisible = this.visible;
		this.res = res;
		if (!wasVisible) {
			return;
		}
		this.hide();
	}
};
