// 3D world scoreboard panel (decui UI.Window + Transform3D).

class WorldScoreboardComponent extends Component {
	className = "WorldScoreboardComponent";

	side = null;
	rows = null;
	res = null;
	colsLabel = "Player          Pts   Mrk   Win";

	static ESCORT_POS = Vector(493.5, 26.0, 13.8);
	static ESCORT_ROT = Vector(-1.5711, 0, 5.35);
	static ESCORT_SIZE = Vector(4.2, 4.8, 0);
	static DEFEND_POS = Vector(506.5, 14.0, 13.8);
	static DEFEND_ROT = Vector(-1.5711, 0, 3.85);
	static DEFEND_SIZE = Vector(4.2, 4.8, 0);

	constructor(o) {
		this.id = o.id;
		this.side = o.side;
		this.rows = o.rows ? o.rows : [];
		this.res = o.res;
		if (o.rawin("colsLabel") && o.colsLabel != null) {
			this.colsLabel = o.colsLabel;
		}
		base.constructor(this.id, o);
		this.metadata.list = "windows";
		this.build();
	}

	function teamColour() {
		return side == "escort" ? SafariTheme.ESCORT : SafariTheme.DEFEND;
	}

	function teamPanelColour() {
		return side == "escort" ? SafariTheme.ESCORT_PANEL : SafariTheme.DEFEND_PANEL;
	}

	function teamRowColour() {
		return side == "escort" ? SafariTheme.ESCORT_ROW : SafariTheme.DEFEND_ROW;
	}

	function teamLabel() {
		return side == "escort" ? "YELLOW team top 10" : "PINK team top 10";
	}

	function transform3D() {
		if (side == "escort") {
			return {
				Position3D = ESCORT_POS,
				Rotation3D = ESCORT_ROT,
				Size3D = ESCORT_SIZE,
			};
		}
		return {
			Position3D = DEFEND_POS,
			Rotation3D = DEFEND_ROT,
			Size3D = DEFEND_SIZE,
		};
	}

	function lineLabel(lineId, y, text, colour = SafariTheme.TEXT, fontSize = 12) {
		return UI.Label({
			id = lineId,
			Position = VectorScreen(8, y),
			Size = VectorScreen(520, 20),
			Text = text,
			TextColour = colour,
			FontSize = fontSize,
			FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD,
		});
	}

	function buildChildren() {
		local c = teamColour();
		local rowC = teamRowColour();
		local children = [
			lineLabel(id + "::title", 20, teamLabel(), c, 13),
			lineLabel(id + "::cols", 45, colsLabel, SafariTheme.TEXT, 11),
		];

		local rowY = 72;
		local rank = 1;
		if (rows.len() == 0) {
			children.append(lineLabel(id + "::empty", rowY, "No records yet", SafariTheme.MUTED, 11));
			return children;
		}

		foreach (row in rows) {
			if (rank > 10) {
				break;
			}
			local name = row.name;
			if (name.len() > 14) {
				name = name.slice(0, 14);
			}
			local line = (rank < 10 ? " " + rank : rank.tostring()) + "  " + name;
			while (line.len() < 22) {
				line += " ";
			}
			line += row.points + "  " + row.marks + "  " + row.wins;
			children.append(lineLabel(id + "::row" + rank, rowY, line, rowC, 11));
			rowY += 24;
			rank++;
		}
		return children;
	}

	function build() {
		UI.Window({
			id = id,
			flags = GUI_FLAG_3D_ENTITY | GUI_FLAG_VISIBLE | GUI_FLAG_TEXT_TAGS,
			Position = VectorScreen(75, 75),
			Size = VectorScreen(540, 420),
			Colour = teamPanelColour(),
			FontName = "Verdana",
			Transform3D = transform3D(),
			children = buildChildren(),
		});
	}

	function destroy() {
		if (UI.findById(id) != null) {
			UI.DeleteByID(id);
		}
	}
}
