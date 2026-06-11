class RoundScoreboardWindow {
	static COL_ESCORT = Colour(110, 210, 255);
	static COL_DEFEND = Colour(255, 110, 110);
	static COL_HEADER = Colour(40, 40, 40);
	static COL_WHITE = Colour(255, 255, 255);

	canvas = null;
	res = null;
	visible = false;

	constructor(res) {
		this.res = res;
	}

	function hide() {
		if (this.canvas != null) {
			this.canvas.Detach();
			this.canvas = null;
		}
		this.visible = false;
		GUI.SetMouseEnabled(false);
	}

	function show(winnerTeam, escortScore, defendScore, reason, players) {
		this.hide();

		local w = this.res.X;
		local h = this.res.Y;
		this.canvas = GUICanvas();
		this.canvas.Position = VectorScreen(0, 0);
		this.canvas.Size = VectorScreen(w, h);

		local backdrop = GUILabel(VectorScreen(0, 0), Colour(0, 0, 0), "");
		backdrop.Size = VectorScreen(w, h);
		this.canvas.AddChild(backdrop);

		local winnerName = winnerTeam == Teams.ESCORT ? "ESCORT" : "DEFENDERS";
		local winnerColour = winnerTeam == Teams.ESCORT ? COL_ESCORT : COL_DEFEND;
		local banner = GUILabel(VectorScreen(0, floor(h * 0.08)), winnerColour,
			winnerName + " WIN — " + escortScore + " vs " + defendScore);
		banner.FontSize = 24;
		banner.Size = VectorScreen(w, 36);
		banner.TextAlignment = GUI_ALIGN_CENTERH;
		banner.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(banner);

		local sub = GUILabel(VectorScreen(0, floor(h * 0.08) + 34), COL_WHITE, reason);
		sub.FontSize = 14;
		sub.Size = VectorScreen(w, 22);
		sub.TextAlignment = GUI_ALIGN_CENTERH;
		sub.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(sub);

		this.drawTeamTable(Teams.ESCORT, escortScore, floor(h * 0.18), players);
		this.drawTeamTable(Teams.DEFEND, defendScore, floor(h * 0.52), players);

		local hint = GUILabel(VectorScreen(w - 120, h - 28), COL_WHITE, "P Close");
		hint.FontSize = 12;
		hint.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(hint);

		this.visible = true;
	}

	function drawTeamTable(team, teamScore, topY, players) {
		local w = this.res.X;
		local tableW = floor(w * 0.78);
		local leftX = floor((w - tableW) / 2);
		local teamColour = team == Teams.ESCORT ? COL_ESCORT : COL_DEFEND;
		local teamLabel = team == Teams.ESCORT ? "ESCORT TEAM" : "DEFEND TEAM";

		local scoreBadge = GUILabel(VectorScreen(leftX - 56, topY + 20), teamColour, teamScore.tostring());
		scoreBadge.FontSize = 28;
		scoreBadge.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(scoreBadge);

		local header = GUILabel(VectorScreen(leftX, topY), teamColour, teamLabel);
		header.FontSize = 16;
		header.Size = VectorScreen(tableW, 24);
		header.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(header);

		local cols = GUILabel(VectorScreen(leftX + 8, topY + 26), COL_WHITE,
			"Player                          Points   Kills   Deaths");
		cols.FontSize = 12;
		cols.Size = VectorScreen(tableW, 18);
		cols.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(cols);

		local rowY = topY + 48;
		local rowNum = 0;
		foreach (p in players) {
			if (p.team != team) {
				continue;
			}
			local name = p.name;
			if (name.len() > 22) {
				name = name.slice(0, 22);
			}
			local line = name + "    " + p.points + "    " + p.kills + "    " + p.deaths;
			local row = GUILabel(VectorScreen(leftX + 8, rowY + rowNum * 20), teamColour, line);
			row.FontSize = 12;
			row.FontFlags = GUI_FFLAG_OUTLINE;
			row.Size = VectorScreen(tableW, 18);
			this.canvas.AddChild(row);
			rowNum++;
			if (rowNum >= 8) {
				break;
			}
		}
		if (rowNum == 0) {
			local empty = GUILabel(VectorScreen(leftX + 8, rowY), COL_WHITE, "(no players)");
			empty.FontSize = 12;
			this.canvas.AddChild(empty);
		}
	}

	function onResize(res) {
		this.res = res;
	}
}
