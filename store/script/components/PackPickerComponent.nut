class PackPickerComponent extends Component {
	className = "PackPickerComponent";

	team = 0;
	currentPack = 1;
	res = null;
	rootId = null;

	constructor(o) {
		this.id = o.id;
		this.rootId = o.id + "::overlay";
		this.team = o.team;
		this.currentPack = o.currentPack;
		this.res = o.res;
		base.constructor(this.id, o);
		this.metadata.list = "canvas";
		this.build();
	}

	function accent() {
		return team == Teams.ESCORT ? SafariTheme.ESCORT : SafariTheme.DEFEND;
	}

	function buildPackCard(packNum, info, x, y, cardW, cardH, isSelected) {
		local cardId = id + "::card" + packNum;
		local ac = accent();
		local cardColour = isSelected ? Colour(ac.R, ac.G, ac.B, 55) : SafariTheme.PANEL_ALT;
		local borderColour = isSelected ? ac : SafariTheme.BORDER;

		local selectFn = function() {
			local s = Stream();
			s.WriteInt(Packets.SELECT_PACK);
			s.WriteInt(packNum);
			Server.SendData(s);
		};

		local card = UI.Canvas({
			id = cardId,
			Position = VectorScreen(x, y),
			Size = VectorScreen(cardW, cardH),
			Colour = cardColour,
			onClick = selectFn,
			children = [
				UI.Label({
					id = cardId + "::title",
					Position = VectorScreen(12, 12),
					Size = VectorScreen(cardW - 24, 22),
					Text = info.title,
					TextColour = ac,
					FontSize = 15,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = cardId + "::weapons",
					Position = VectorScreen(12, 40),
					Size = VectorScreen(cardW - 24, cardH - 90),
					Text = info.weapons,
					TextColour = SafariTheme.MUTED,
					FontSize = 12,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
				UI.Button({
					id = cardId + "::btn",
					Text = "pack" + packNum,
					TextColour = SafariTheme.TEXT,
					Position = VectorScreen(12, cardH - 44),
					Size = VectorScreen(cardW - 24, 32),
					Colour = isSelected ? SafariTheme.SELECTED : SafariTheme.BUTTON,
					onClick = selectFn,
					onHoverOver = function() {
						if (!isSelected) {
							this.Colour = SafariTheme.BUTTON_HOVER;
						}
					},
					onHoverOut = function() {
						if (!isSelected) {
							this.Colour = SafariTheme.BUTTON;
						}
					},
				}),
			],
		});
		card.addBorders({ size = 2, color = borderColour });
		return card;
	}

	function build() {
		local panelW = 720;
		local panelH = 400;
		local packs = packListForTeam(team);
		local cardW = 210;
		local gap = 18;
		local startX = (panelW - (cardW * 3 + gap * 2)) / 2;

		local cardChildren = [];
		for (local i = 0; i < 3; i++) {
			local packNum = i + 1;
			cardChildren.append(buildPackCard(packNum, packs[i], startX + i * (cardW + gap), 118, cardW, 230, packNum == currentPack));
		}

		local panelChildren = [
			UI.Label({
				id = id + "::title",
				align = "center",
				Position = VectorScreen(0, 22),
				Size = VectorScreen(panelW, 30),
				Text = "SELECT YOUR WEAPON PACK",
				TextColour = SafariTheme.TEXT,
				FontSize = 22,
				FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
			}),
			UI.Label({
				id = id + "::status",
				align = "center",
				Position = VectorScreen(40, 56),
				Size = VectorScreen(panelW - 80, 22),
				Text = "",
				TextColour = SafariTheme.STATUS,
				FontSize = 13,
				FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD,
			}),
		];
		foreach (c in cardChildren) {
			panelChildren.append(c);
		}
		panelChildren.append(
			UI.Label({
				id = id + "::hint",
				Position = VectorScreen(panelW - 110, panelH - 28),
				Size = VectorScreen(100, 18),
				Text = "[P] Close",
				TextColour = SafariTheme.MUTED,
				FontSize = 12,
				FontFlags = GUI_FFLAG_OUTLINE,
			})
		);

		local panel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(panelW, panelH),
			Colour = SafariTheme.PANEL,
			children = panelChildren,
		});
		panel.addBorders({ size = 2, color = accent() });

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
