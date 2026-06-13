class LeaderboardComponent extends Component {
	className = "LeaderboardComponent";

	escortRows = null;
	defendRows = null;
	res = null;
	rootId = null;

	// Column widths (monospace-style padding)
	RANK_W = 5;
	NAME_W = 28;
	PTS_W = 10;
	MRK_W = 9;
	WIN_W = 7;

	constructor(o) {
		this.id = o.id;
		this.rootId = o.id + "::overlay";
		this.escortRows = o.escortRows;
		this.defendRows = o.defendRows;
		this.res = o.res;
		base.constructor(this.id, o);
		this.metadata.list = "canvas";
		this.build();
	}

	function teamTotal(rows) {
		local total = 0;
		foreach (row in rows) {
			total += row.points;
		}
		return total;
	}

	function padRight(text, width) {
		local s = text;
		while (s.len() < width) {
			s += " ";
		}
		return s;
	}

	function padLeft(text, width) {
		local s = text.tostring();
		while (s.len() < width) {
			s = " " + s;
		}
		return s;
	}

	function formatName(name) {
		if (name.len() > NAME_W - 1) {
			return name.slice(0, NAME_W - 1);
		}
		return name;
	}

	function formatHeader() {
		return padRight("", RANK_W)
			+ padRight("Name", NAME_W)
			+ padLeft("Points", PTS_W)
			+ padLeft("Mrk", MRK_W)
			+ padLeft("Win", WIN_W);
	}

	function formatRow(rank, row) {
		return padLeft(rank.tostring() + ".", RANK_W)
			+ padRight(formatName(row.name), NAME_W)
			+ padLeft(row.points, PTS_W)
			+ padLeft(row.marks, MRK_W)
			+ padLeft(row.wins, WIN_W);
	}

	function buildTeamTable(side, rows, panelId, panelY, tableW, scoreX) {
		local isPink = side == "defend";
		local panelColour = isPink ? SafariTheme.DEFEND_PANEL : SafariTheme.ESCORT_PANEL;
		local textColour = isPink ? SafariTheme.DEFEND_ROW : SafariTheme.ESCORT_ROW;
		local accent = isPink ? SafariTheme.DEFEND : SafariTheme.ESCORT;
		local total = teamTotal(rows);
		local rowCount = rows.len() > 10 ? 10 : (rows.len() == 0 ? 1 : rows.len());
		local scoreW = 72;
		local contentX = scoreW + 10;
		local contentW = tableW - contentX - 10;
		local headerY = 10;
		local rowY = 32;
		local tableH = rowY + rowCount * 22 + 10;

		local children = [
			UI.Label({
				id = panelId + "::score",
				Position = VectorScreen(8, 14),
				Size = VectorScreen(scoreW - 8, 48),
				Text = total.tostring(),
				TextColour = accent,
				FontSize = 34,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
			UI.Label({
				id = panelId + "::cols",
				Position = VectorScreen(contentX, headerY),
				Size = VectorScreen(contentW, 16),
				Text = formatHeader(),
				TextColour = SafariTheme.TEXT,
				FontSize = 11,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
		];

		if (rows.len() == 0) {
			children.append(UI.Label({
				id = panelId + "::empty",
				Position = VectorScreen(contentX, rowY),
				Size = VectorScreen(contentW, 18),
				Text = "No records yet",
				TextColour = SafariTheme.MUTED,
				FontSize = 11,
				FontFlags = GUI_FFLAG_OUTLINE,
			}));
		} else {
			local rank = 1;
			local rowPos = rowY;
			foreach (row in rows) {
				if (rank > 10) {
					break;
				}
				children.append(UI.Label({
					id = panelId + "::row" + rank,
					Position = VectorScreen(contentX, rowPos),
					Size = VectorScreen(contentW, 20),
					Text = formatRow(rank, row),
					TextColour = textColour,
					FontSize = 11,
					FontFlags = GUI_FFLAG_OUTLINE,
				}));
				rowPos += 22;
				rank++;
			}
		}

		return UI.Canvas({
			id = panelId,
			Position = VectorScreen(scoreX, panelY),
			Size = VectorScreen(tableW, tableH),
			Colour = panelColour,
			children = children,
		});
	}

	function build() {
		local tableW = 680;
		local scoreX = 0;
		local defendRowsCount = defendRows.len() > 10 ? 10 : (defendRows.len() == 0 ? 1 : defendRows.len());
		local escortRowsCount = escortRows.len() > 10 ? 10 : (escortRows.len() == 0 ? 1 : escortRows.len());
		local defendH = 32 + defendRowsCount * 22 + 10;
		local escortH = 32 + escortRowsCount * 22 + 10;
		local gap = 12;
		local headerH = 52;
		local mainW = tableW;
		local mainH = headerH + defendH + gap + escortH + 8;

		local defendTable = buildTeamTable("defend", defendRows, id + "::defend", headerH, tableW, scoreX);
		local escortTable = buildTeamTable("escort", escortRows, id + "::escort", headerH + defendH + gap, tableW, scoreX);

		local mainPanel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(mainW, mainH),
			Colour = Colour(0, 0, 0, 0),
			children = [
				UI.Label({
					id = id + "::heading",
					Position = VectorScreen(0, 6),
					Size = VectorScreen(mainW, 26),
					Text = "LEADERBOARDS",
					TextColour = SafariTheme.TEXT,
					FontSize = 22,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::hint",
					Position = VectorScreen(0, 32),
					Size = VectorScreen(mainW, 16),
					Text = "/leaderboard to close  |  P to close",
					TextColour = SafariTheme.MUTED,
					FontSize = 11,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
				defendTable,
				escortTable,
			],
		});

		UI.Canvas({
			id = rootId,
			align = "center",
			Size = VectorScreen(res.X, res.Y),
			Position = VectorScreen(0, 0),
			Colour = SafariTheme.OVERLAY,
			children = [mainPanel],
		});

		GUI.SetMouseEnabled(true);
	}

	function destroy() {
		if (UI.findById(rootId) != null) {
			UI.DeleteByID(rootId);
		}
	}
}
