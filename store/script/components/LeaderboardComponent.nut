class LeaderboardComponent extends Component {
	className = "LeaderboardComponent";

	escortRows = null;
	defendRows = null;
	res = null;
	rootId = null;

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
		if (name.len() > 18) {
			return name.slice(0, 18);
		}
		return name;
	}

	function formatRow(rank, row) {
		local rankStr = rank < 10 ? rank.tostring() + "." : rank.tostring() + ".";
		return padRight(rankStr + " " + formatName(row.name), 24)
			+ padLeft(row.points, 7)
			+ padLeft(row.marks, 6)
			+ padLeft(row.wins, 5);
	}

	function buildTeamTable(side, rows, panelId, y, tableW, scoreX) {
		local isPink = side == "defend";
		local panelColour = isPink ? SafariTheme.DEFEND_PANEL : SafariTheme.ESCORT_PANEL;
		local textColour = isPink ? SafariTheme.DEFEND_ROW : SafariTheme.ESCORT_ROW;
		local accent = isPink ? SafariTheme.DEFEND : SafariTheme.ESCORT;
		local label = isPink ? "PINK" : "YELLOW";
		local total = teamTotal(rows);
		local rowCount = rows.len() > 10 ? 10 : (rows.len() == 0 ? 1 : rows.len());
		local tableH = 52 + rowCount * 22 + 12;
		local scoreW = 72;

		local children = [
			UI.Label({
				id = panelId + "::score",
				Position = VectorScreen(8, 18),
				Size = VectorScreen(scoreW - 8, 48),
				Text = total.tostring(),
				TextColour = accent,
				FontSize = 34,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
			UI.Label({
				id = panelId + "::team",
				Position = VectorScreen(scoreW + 8, 10),
				Size = VectorScreen(tableW - scoreW - 16, 24),
				Text = label + " TEAM",
				TextColour = accent,
				FontSize = 18,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
			UI.Label({
				id = panelId + "::cols",
				Position = VectorScreen(scoreW + 8, 34),
				Size = VectorScreen(tableW - scoreW - 16, 16),
				Text = padRight("Name", 24) + padLeft("Points", 7) + padLeft("Mrk", 6) + padLeft("Win", 5),
				TextColour = SafariTheme.TEXT,
				FontSize = 11,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
		];

		local rowY = 54;
		if (rows.len() == 0) {
			children.append(UI.Label({
				id = panelId + "::empty",
				Position = VectorScreen(scoreW + 8, rowY),
				Size = VectorScreen(tableW - scoreW - 16, 18),
				Text = "No records yet",
				TextColour = SafariTheme.MUTED,
				FontSize = 11,
				FontFlags = GUI_FFLAG_OUTLINE,
			}));
		} else {
			local rank = 1;
			foreach (row in rows) {
				if (rank > 10) {
					break;
				}
				children.append(UI.Label({
					id = panelId + "::row" + rank,
					Position = VectorScreen(scoreW + 8, rowY),
					Size = VectorScreen(tableW - scoreW - 16, 20),
					Text = formatRow(rank, row),
					TextColour = textColour,
					FontSize = 11,
					FontFlags = GUI_FFLAG_OUTLINE,
				}));
				rowY += 22;
				rank++;
			}
		}

		return UI.Canvas({
			id = panelId,
			Position = VectorScreen(scoreX, y),
			Size = VectorScreen(tableW, tableH),
			Colour = panelColour,
			children = children,
		});
	}

	function build() {
		local tableW = 640;
		local scoreX = 16;
		local defendRowsCount = defendRows.len() > 10 ? 10 : (defendRows.len() == 0 ? 1 : defendRows.len());
		local escortRowsCount = escortRows.len() > 10 ? 10 : (escortRows.len() == 0 ? 1 : escortRows.len());
		local defendH = 52 + defendRowsCount * 22 + 12;
		local escortH = 52 + escortRowsCount * 22 + 12;
		local gap = 10;
		local mainW = tableW + 32;
		local mainH = 88 + defendH + gap + escortH;

		// Flag Raids order: pink (defend) on top, yellow (escort) below.
		local defendTable = buildTeamTable("defend", defendRows, id + "::defend", 72, tableW, scoreX);
		local escortTable = buildTeamTable("escort", escortRows, id + "::escort", 72 + defendH + gap, tableW, scoreX);

		local mainPanel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(mainW, mainH),
			Colour = Colour(0, 0, 0, 0),
			children = [
				UI.Label({
					id = id + "::heading",
					align = "center",
					Position = VectorScreen(0, 12),
					Size = VectorScreen(mainW, 28),
					Text = "LEADERBOARDS",
					TextColour = SafariTheme.TEXT,
					FontSize = 22,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::hint",
					align = "center",
					Position = VectorScreen(0, 42),
					Size = VectorScreen(mainW, 18),
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
