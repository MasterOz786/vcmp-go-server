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

	function formatRow(rank, row) {
		local name = row.name;
		if (name.len() > 16) {
			name = name.slice(0, 16);
		}
		local rankStr = rank < 10 ? " " + rank : rank.tostring();
		local line = rankStr + "  " + name;
		while (line.len() < 24) {
			line += " ";
		}
		return line + row.points + "  " + row.marks + "  " + row.wins;
	}

	function buildTeamPanel(side, rows, panelId, x, y, pw, ph) {
		local colour = side == "escort" ? SafariTheme.ESCORT : SafariTheme.DEFEND;
		local label = side == "escort" ? "ESCORT" : "DEFEND";
		local children = [
			UI.Label({
				id = panelId + "::stripe",
				Position = VectorScreen(0, 0),
				Size = VectorScreen(pw, 3),
				Text = "",
				Colour = colour,
			}),
			UI.Label({
				id = panelId + "::title",
				Position = VectorScreen(12, 12),
				Size = VectorScreen(pw - 24, 22),
				Text = label,
				TextColour = colour,
				FontSize = 16,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
			UI.Label({
				id = panelId + "::cols",
				Position = VectorScreen(12, 36),
				Size = VectorScreen(pw - 24, 16),
				Text = "#  Player              Pts  Mrk  Win",
				TextColour = SafariTheme.MUTED,
				FontSize = 10,
				FontFlags = GUI_FFLAG_OUTLINE,
			}),
		];

		local rowY = 56;
		if (rows.len() == 0) {
			children.append(UI.Label({
				id = panelId + "::empty",
				Position = VectorScreen(12, rowY),
				Size = VectorScreen(pw - 24, 18),
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
					Position = VectorScreen(12, rowY),
					Size = VectorScreen(pw - 24, 18),
					Text = formatRow(rank, row),
					TextColour = colour,
					FontSize = 11,
					FontFlags = GUI_FFLAG_OUTLINE,
				}));
				rowY += 20;
				rank++;
			}
		}

		local panel = UI.Canvas({
			id = panelId,
			Position = VectorScreen(x, y),
			Size = VectorScreen(pw, ph),
			Colour = SafariTheme.PANEL_ALT,
			children = children,
		});
		panel.addBorders({ size = 2, color = colour });
		return panel;
	}

	function build() {
		local panelW = 320;
		local panelH = 280;
		local mainW = panelW * 2 + 48;
		local mainH = panelH + 96;

		local escortPanel = buildTeamPanel("escort", escortRows, id + "::escort", 16, 72, panelW, panelH);
		local defendPanel = buildTeamPanel("defend", defendRows, id + "::defend", 16 + panelW + 16, 72, panelW, panelH);

		local mainPanel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(mainW, mainH),
			Colour = SafariTheme.PANEL,
			children = [
				UI.Label({
					id = id + "::heading",
					align = "center",
					Position = VectorScreen(0, 16),
					Size = VectorScreen(mainW, 28),
					Text = "SAFARI LEADERBOARDS",
					TextColour = SafariTheme.TEXT,
					FontSize = 22,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::hint",
					align = "center",
					Position = VectorScreen(0, 44),
					Size = VectorScreen(mainW, 18),
					Text = "/leaderboard to close  |  P to close",
					TextColour = SafariTheme.MUTED,
					FontSize = 11,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
				escortPanel,
				defendPanel,
			],
		});
		mainPanel.addBorders({ size = 2, color = SafariTheme.BORDER });

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
