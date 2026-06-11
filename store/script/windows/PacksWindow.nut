class PacksWindow {
	static COL_SELECTED = Colour(240, 175, 238);
	static PACK_COUNT = 3;
	static CANVAS_W = 370;
	static CANVAS_H = 200;

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

		this.statusLabel = GUILabel(VectorScreen(0, 93), COL_SELECTED, "");
		this.statusLabel.FontSize = 16;
		this.statusLabel.Size = VectorScreen(CANVAS_W - 67, 200);
		this.statusLabel.TextAlignment = GUI_ALIGN_CENTERH;
		this.statusLabel.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;

		this.canvas = GUICanvas();
		this.canvas.Position = VectorScreen((res.X / 2) - (CANVAS_W / 2), (res.Y / 2) - (100 / 2) + 30);
		this.canvas.Size = VectorScreen(CANVAS_W, CANVAS_H);
		this.canvas.AddChild(this.statusLabel);

		local hSpace = 10;
		for (local i = 1; i <= PACK_COUNT; i++) {
			local btn = GUIButton(VectorScreen(hSpace, 10), VectorScreen(45, 26), COL_SELECTED, "pack" + i);
			btn.TextColour = Colour(0, 0, 0);
			btn.FontName = "Verdana";
			btn.FontFlags = GUI_FFLAG_BOLD;
			if (i == currentPack) {
				btn.Text = "pack" + i + " *";
			}
			hSpace += 51;
			this.canvas.AddChild(btn);
			this.packButtons.append(btn);
		}

		GUI.SetMouseEnabled(true);
	}

	function updatePositions(res) {
		this.res = res;
		if (this.canvas != null) {
			this.canvas.Position = VectorScreen((res.X / 2) - (CANVAS_W / 2), (res.Y / 2) - (100 / 2) + 30);
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
