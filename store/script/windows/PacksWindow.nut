class PacksWindow {
	static COL_HEADER = Colour(255, 140, 0);
	static COL_BG = Colour(10, 10, 10);
	static COL_TITLE = Colour(255, 255, 255);
	static COL_WEAPONS = Colour(230, 230, 230);
	static COL_SELECT = Colour(240, 175, 238);

	canvas = null;
	res = null;
	team = 0;
	currentPack = 1;
	packButtons = [];
	statusLabel = null;

	constructor(res) {
		this.res = res;
	}

	function clear() {
		this.canvas = null;
		this.packButtons.clear();
		this.statusLabel = null;
		GUI.SetMouseEnabled(false);
	}

	function setStatus(msg) {
		if (this.statusLabel != null) {
			this.statusLabel.Text = msg;
		}
	}

	function createWindow(team, currentPack) {
		if (this.canvas != null) {
			return;
		}
		this.team = team;
		this.currentPack = currentPack;

		local packW = 200;
		local packH = 170;
		local gap = 16;
		local totalW = packW * 3 + gap * 2;
		local startX = (this.res.X / 2) - (totalW / 2);
		local startY = (this.res.Y / 2) - 90;

		this.canvas = GUICanvas();
		this.canvas.Position = VectorScreen(startX - 20, startY - 60);
		this.canvas.Size = VectorScreen(totalW + 40, packH + 110);

		local title = GUILabel(VectorScreen(20, 0), COL_TITLE, "SELECT YOUR WEAPON PACK");
		title.FontSize = 22;
		title.Size = VectorScreen(totalW, 30);
		title.TextAlignment = GUI_ALIGN_CENTERH;
		title.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;
		this.canvas.AddChild(title);

		this.statusLabel = GUILabel(VectorScreen(20, 28), Colour(255, 200, 120), "");
		this.statusLabel.FontSize = 13;
		this.statusLabel.Size = VectorScreen(totalW, 20);
		this.statusLabel.TextAlignment = GUI_ALIGN_CENTERH;
		this.canvas.AddChild(this.statusLabel);

		local closeHint = GUILabel(VectorScreen(totalW - 70, packH + 72), COL_TITLE, "P Close");
		closeHint.FontSize = 12;
		closeHint.FontFlags = GUI_FFLAG_OUTLINE;
		this.canvas.AddChild(closeHint);

		local packs = packListForTeam(team);
		for (local i = 0; i < packs.len(); i++) {
			local px = 20 + i * (packW + gap);
			local py = 52;
			local packNum = i + 1;
			local info = packs[i];

			local panel = GUICanvas();
			panel.Position = VectorScreen(px, py);
			panel.Size = VectorScreen(packW, packH);

			local header = GUILabel(VectorScreen(0, 0), COL_HEADER, info.title);
			header.Size = VectorScreen(packW, 26);
			header.TextAlignment = GUI_ALIGN_CENTERH | GUI_ALIGN_CENTERV;
			header.FontFlags = GUI_FFLAG_BOLD;
			header.FontSize = 14;
			panel.AddChild(header);

			local weapons = GUILabel(VectorScreen(8, 34), COL_WEAPONS, info.weapons);
			weapons.Size = VectorScreen(packW - 16, 80);
			weapons.FontSize = 12;
			weapons.FontFlags = GUI_FFLAG_OUTLINE;
			weapons.AddFlags(GUI_FLAG_WRAP);
			panel.AddChild(weapons);

			local btn = GUIButton(VectorScreen(40, packH - 34), VectorScreen(120, 26), COL_SELECT, "pack" + packNum);
			btn.TextColour = Colour(0, 0, 0);
			btn.FontFlags = GUI_FFLAG_BOLD;
			if (packNum == currentPack) {
				btn.Text = "pack" + packNum + " *";
			}
			panel.AddChild(btn);
			this.packButtons.append(btn);

			this.canvas.AddChild(panel);
		}

		GUI.SetMouseEnabled(true);
	}

	function updatePositions(res) {
		this.res = res;
		if (this.canvas != null) {
			this.clear();
		}
	}

	function requestToggle() {
		local s = Stream();
		s.WriteInt(Packets.REQUEST_SHOW_PACKS);
		Server.SendData(s);
	}

	function selectPack(packNum) {
		local s = Stream();
		s.WriteInt(Packets.SELECT_PACK);
		s.WriteInt(packNum);
		Server.SendData(s);
	}
}
