// 2D overlay for /leaderboard. 3D boards live in RoundScoreboardController (CTF pattern).

class LobbyLeaderboardController {
	static COL_ESCORT = Colour(110, 210, 255);
	static COL_DEFEND = Colour(255, 110, 110);
	static COL_WHITE = Colour(255, 255, 255);
	static COL_MUTED = Colour(155, 168, 188);
	static COL_BG = Colour(14, 18, 28);
	static COL_PANEL = Colour(22, 28, 42);

	canvas = null;
	res = null;
	visible = false;

	constructor(res) {
		this.res = res;
	}

	function hideOverlay() {
		this.canvas = null;
		GUI.SetMouseEnabled(false);
		this.visible = false;
	}

	function hide() {
		this.hideOverlay();
	}

	function showOverlay(escortRows, defendRows) {
		this.hideOverlay();

		local w = this.res.X;
		local h = this.res.Y;
		local panelW = floor(w * 0.38);
		local panelH = floor(h * 0.55);
		local topY = floor(h * 0.14);
		local leftX = floor(w * 0.08);
		local rightX = w - panelW - leftX;

		this.canvas = GUICanvas();
		this.canvas.Position = VectorScreen(0, 0);
		this.canvas.Size = VectorScreen(w, h);

		local backdrop = GUILabel(VectorScreen(0, 0), COL_BG, "");
		backdrop.Size = VectorScreen(w, h);
		backdrop.Alpha = 200;
		this.canvas.AddChild(backdrop);

		local title = GUILabel(VectorScreen(0, floor(h * 0.04)), COL_WHITE, "SAFARI LEADERBOARDS");
		title.FontSize = 22;
		title.Size = VectorScreen(w, 30);
		title.TextAlignment = GUI_ALIGN_CENTERH;
		title.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(title);

		local hint = GUILabel(VectorScreen(0, floor(h * 0.04) + 30), COL_MUTED,
			"All-time stats  |  /leaderboard to close  |  P to close");
		hint.FontSize = 11;
		hint.Size = VectorScreen(w, 18);
		hint.TextAlignment = GUI_ALIGN_CENTERH;
		hint.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(hint);

		this.drawPanel(leftX, topY, panelW, panelH, "ESCORT", COL_ESCORT, escortRows);
		this.drawPanel(rightX, topY, panelW, panelH, "DEFEND", COL_DEFEND, defendRows);

		GUI.SetMouseEnabled(true);
		this.visible = true;
	}

	function drawPanel(x, y, pw, ph, label, colour, rows) {
		local panel = GUILabel(VectorScreen(x, y), COL_PANEL, "");
		panel.Size = VectorScreen(pw, ph);
		panel.Alpha = 220;
		this.canvas.AddChild(panel);

		local stripe = GUILabel(VectorScreen(x, y), colour, "");
		stripe.Size = VectorScreen(pw, 3);
		this.canvas.AddChild(stripe);

		local header = GUILabel(VectorScreen(x + 12, y + 10), colour, label);
		header.FontSize = 16;
		header.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(header);

		local cols = GUILabel(VectorScreen(x + 12, y + 34), COL_MUTED,
			"#  Player              Pts  Mrk  Win");
		cols.FontSize = 10;
		cols.Size = VectorScreen(pw - 20, 14);
		cols.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(cols);

		local rowY = y + 54;
		if (rows.len() == 0) {
			local empty = GUILabel(VectorScreen(x + 12, rowY), COL_MUTED, "No records yet");
			empty.FontSize = 11;
			this.canvas.AddChild(empty);
			return;
		}

		local rank = 1;
		foreach (row in rows) {
			if (rank > 10) {
				break;
			}
			local name = row.name;
			if (name.len() > 16) {
				name = name.slice(0, 16);
			}
			local rankStr = rank < 10 ? " " + rank : rank.tostring();
			local line = rankStr + "  " + name;
			while (line.len() < 24) {
				line += " ";
			}
			line += row.points + "  " + row.marks + "  " + row.wins;
			local rowLabel = GUILabel(VectorScreen(x + 12, rowY), colour, line);
			rowLabel.FontSize = 11;
			rowLabel.FontFlags = GUI_FFLAG_OUTLINE;
			rowLabel.Size = VectorScreen(pw - 20, 16);
			this.canvas.AddChild(rowLabel);
			rowY += 20;
			rank++;
		}
	}

	function onResize(res) {
		this.res = res;
	}
}
