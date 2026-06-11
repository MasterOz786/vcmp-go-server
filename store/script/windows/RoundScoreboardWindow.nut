class TeamRoundScoreboardWindow {
	static COL_TITLEBAR = Colour(68, 68, 68);
	static COL_ESCORT = Colour(110, 210, 255);
	static COL_DEFEND = Colour(255, 110, 110);
	static COL_WHITE = Colour(255, 255, 255);

	// patrol_default.json spawn area (not CTF map coords).
	static ESCORT_ROT = Vector(-1.5711, 0, 4.6);
	static ESCORT_SIZE = Vector(3.6, 4, 0);
	static ESCORT_POS = Vector(-1145.0, 15.0, 16.0);

	static DEFEND_ROT = Vector(-1.5711, 0, 3.1);
	static DEFEND_SIZE = Vector(4, 5.0, 0);
	static DEFEND_POS = Vector(-920.0, 175.0, 24.0);

	window = null;
	elements = [];
	players = [];
	res = null;
	side = null;
	rowY = 75;

	constructor(res, side) {
		this.res = res;
		this.side = side;
	}

	function teamColour() {
		return side == "escort" ? COL_ESCORT : COL_DEFEND;
	}

	function teamTitle() {
		return side == "escort" ? "Escort team scoreboard" : "Defend team scoreboard";
	}

	function matchesTeam(team) {
		if (side == "escort") {
			return team == Teams.ESCORT;
		}
		return team == Teams.DEFEND;
	}

	function clearWindow() {
		window = null;
		elements.clear();
		players.clear();
		rowY = 75;
	}

	function createWindow() {
		if (window != null) {
			return;
		}

		window = GUIWindow(VectorScreen(75, 75), VectorScreen(540, 420), COL_TITLEBAR, "", GUI_FLAG_TEXT_TAGS);
		window.AddFlags(GUI_FLAG_VISIBLE);
		window.FontName = "Verdana";
		window.FontSize = 10;
		window.FontFlags = GUI_FFLAG_BOLD;
		window.Text = teamTitle();
		window.TitleColour = COL_TITLEBAR;
		window.Alpha = 255;
		window.RemoveFlags(GUI_FLAG_DRAGGABLE | GUI_FLAG_WINDOW_CLOSEBTN | GUI_FLAG_WINDOW_RESIZABLE | GUI_FLAG_WINDOW_TITLEBAR);

		drawHeader();

		window.AddFlags(GUI_FLAG_3D_ENTITY);
		if (side == "escort") {
			window.Set3DTransform(ESCORT_POS, ESCORT_ROT, ESCORT_SIZE);
		} else {
			window.Set3DTransform(DEFEND_POS, DEFEND_ROT, DEFEND_SIZE);
		}
	}

	function addToWindow(el) {
		window.AddChild(el);
		elements.append(el);
	}

	function drawLine(hSpace, vSpace) {
		local line = GUILabel(VectorScreen(hSpace, vSpace + 11), Colour(10, 10, 10),
			"_________________________________________________________________________________________");
		addToWindow(line);
	}

	function applyLabelStyle(label) {
		label.Size = VectorScreen(100, 20);
		label.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
	}

	function drawHeader() {
		drawLine(-5, 20);
		drawLine(-5, 44);

		local title = GUILabel(VectorScreen(3, 20), COL_WHITE, side + " team top 10");
		local nameCol = GUILabel(VectorScreen(11, 45), COL_WHITE, "Player");
		local ptsCol = GUILabel(VectorScreen(146, 45), COL_WHITE, "Points");
		local killsCol = GUILabel(VectorScreen(218, 45), COL_WHITE, "Kills");
		local deathsCol = GUILabel(VectorScreen(294, 45), COL_WHITE, "Deaths");

		nameCol.FontSize = 12;
		ptsCol.FontSize = 12;
		killsCol.FontSize = 12;
		deathsCol.FontSize = 12;

		applyLabelStyle(title);
		applyLabelStyle(nameCol);
		applyLabelStyle(ptsCol);
		applyLabelStyle(killsCol);
		applyLabelStyle(deathsCol);

		addToWindow(title);
		addToWindow(nameCol);
		addToWindow(ptsCol);
		addToWindow(killsCol);
		addToWindow(deathsCol);
	}

	function hasPlayer(name) {
		foreach (p in players) {
			if (p.name == name) {
				return true;
			}
		}
		return false;
	}

	function pushPlayer(name, team, points, kills, deaths) {
		if (!matchesTeam(team)) {
			return;
		}
		if (hasPlayer(name)) {
			return;
		}
		if (players.len() >= 10) {
			return;
		}

		createWindow();

		local stats = { name = name, team = team, points = points, kills = kills, deaths = deaths };
		players.append(stats);

		drawLine(-5, rowY);
		drawRow(rowY, stats);
		rowY += 30;
	}

	function drawRow(vSpace, stats) {
		local c = teamColour();
		local nickname = stats.name;
		if (nickname.len() > 16) {
			nickname = nickname.slice(0, 16);
		}

		local nameLabel = GUILabel(VectorScreen(13, vSpace), c, nickname);
		nameLabel.TextAlignment = GUI_ALIGN_LEFT;
		applyLabelStyle(nameLabel);
		addToWindow(nameLabel);

		local pointsLabel = GUILabel(VectorScreen(150, vSpace), c, stats.points + "");
		applyLabelStyle(pointsLabel);
		addToWindow(pointsLabel);

		local killsLabel = GUILabel(VectorScreen(210, vSpace), c, stats.kills + "");
		applyLabelStyle(killsLabel);
		addToWindow(killsLabel);

		local deathsLabel = GUILabel(VectorScreen(290, vSpace), c, stats.deaths + "");
		applyLabelStyle(deathsLabel);
		addToWindow(deathsLabel);
	}
}

class RoundScoreboardController {
	static COL_ESCORT = Colour(110, 210, 255);
	static COL_DEFEND = Colour(255, 110, 110);
	static COL_WHITE = Colour(255, 255, 255);

	escortBoard = null;
	defendBoard = null;
	bannerCanvas = null;
	res = null;
	visible = false;

	constructor(res) {
		this.res = res;
		this.escortBoard = TeamRoundScoreboardWindow(res, "escort");
		this.defendBoard = TeamRoundScoreboardWindow(res, "defend");
	}

	function hide() {
		this.escortBoard.clearWindow();
		this.defendBoard.clearWindow();
		if (this.bannerCanvas != null) {
			this.bannerCanvas = null;
		}
		this.visible = false;
	}

	function show(winnerTeam, escortScore, defendScore, reason, players) {
		this.hide();

		this.escortBoard.createWindow();
		this.defendBoard.createWindow();

		foreach (p in players) {
			if (p.team == Teams.ESCORT) {
				this.escortBoard.pushPlayer(p.name, p.team, p.points, p.kills, p.deaths);
			} else if (p.team == Teams.DEFEND) {
				this.defendBoard.pushPlayer(p.name, p.team, p.points, p.kills, p.deaths);
			}
		}

		local w = this.res.X;
		local h = this.res.Y;
		this.bannerCanvas = GUICanvas();
		this.bannerCanvas.Position = VectorScreen(0, 0);
		this.bannerCanvas.Size = VectorScreen(w, h);

		local backdrop = GUILabel(VectorScreen(0, 0), Colour(0, 0, 0), "");
		backdrop.Size = VectorScreen(w, h);
		backdrop.Alpha = 180;
		this.bannerCanvas.AddChild(backdrop);

		local winnerName = winnerTeam == Teams.ESCORT ? "ESCORT" : "DEFENDERS";
		local winnerColour = winnerTeam == Teams.ESCORT ? COL_ESCORT : COL_DEFEND;
		local banner = GUILabel(VectorScreen(0, floor(h * 0.06)), winnerColour,
			winnerName + " WIN — " + escortScore + " vs " + defendScore);
		banner.FontSize = 24;
		banner.Size = VectorScreen(w, 36);
		banner.TextAlignment = GUI_ALIGN_CENTERH;
		banner.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.bannerCanvas.AddChild(banner);

		local sub = GUILabel(VectorScreen(0, floor(h * 0.06) + 34), COL_WHITE, reason);
		sub.FontSize = 14;
		sub.Size = VectorScreen(w, 22);
		sub.TextAlignment = GUI_ALIGN_CENTERH;
		sub.FontFlags = GUI_FFLAG_OUTLINE;
		this.bannerCanvas.AddChild(sub);

		this.drawTeamTable(Teams.ESCORT, escortScore, floor(h * 0.16), players);
		this.drawTeamTable(Teams.DEFEND, defendScore, floor(h * 0.50), players);

		local hint = GUILabel(VectorScreen(w - 120, h - 28), COL_WHITE, "P Close");
		hint.FontSize = 12;
		hint.FontFlags = GUI_FFLAG_OUTLINE;
		this.bannerCanvas.AddChild(hint);

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
		this.bannerCanvas.AddChild(scoreBadge);

		local header = GUILabel(VectorScreen(leftX, topY), teamColour, teamLabel);
		header.FontSize = 16;
		header.Size = VectorScreen(tableW, 24);
		header.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.bannerCanvas.AddChild(header);

		local cols = GUILabel(VectorScreen(leftX + 8, topY + 26), COL_WHITE,
			"Player                          Points   Kills   Deaths");
		cols.FontSize = 12;
		cols.Size = VectorScreen(tableW, 18);
		cols.FontFlags = GUI_FFLAG_OUTLINE;
		this.bannerCanvas.AddChild(cols);

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
			this.bannerCanvas.AddChild(row);
			rowNum++;
			if (rowNum >= 8) {
				break;
			}
		}
		if (rowNum == 0) {
			local empty = GUILabel(VectorScreen(leftX + 8, rowY), COL_WHITE, "(no players)");
			empty.FontSize = 12;
			this.bannerCanvas.AddChild(empty);
		}
	}

	function onResize(res) {
		this.res = res;
		this.escortBoard.res = res;
		this.defendBoard.res = res;
	}
}
