class PacksWindow {
	component = null;
	res = null;
	team = 0;
	currentPack = 1;
	packButtons = [];
	rootId = "safari::packpicker";

	constructor(res) {
		this.res = res;
	}

	function clear() {
		if (this.component != null) {
			this.component.destroy();
			this.component = null;
		} else if (UI.findById(this.rootId + "::overlay") != null) {
			UI.DeleteByID(this.rootId + "::overlay");
		}
		this.packButtons.clear();
		GUI.SetMouseEnabled(false);
	}

	function setStatus(msg) {
		local lbl = UI.Label(this.rootId + "::status");
		if (lbl != null) {
			lbl.Text = msg;
		}
	}

	function collectPackButtons() {
		this.packButtons.clear();
		for (local i = 1; i <= 3; i++) {
			local btn = UI.Button(this.rootId + "::card" + i + "::btn");
			if (btn != null) {
				this.packButtons.append(btn);
			}
		}
	}

	function createWindow(team, currentPack) {
		if (this.component != null) {
			return;
		}
		this.team = team;
		this.currentPack = currentPack;
		this.component = PackPickerComponent({
			id = this.rootId,
			team = team,
			currentPack = currentPack,
			res = this.res,
		});
		this.collectPackButtons();
	}

	function updatePositions(res) {
		this.res = res;
		if (this.component == null) {
			return;
		}
		local team = this.team;
		local currentPack = this.currentPack;
		local status = this.setStatusText();
		this.clear();
		this.createWindow(team, currentPack);
		this.setStatus(status);
	}

	function setStatusText() {
		local lbl = UI.Label(this.rootId + "::status");
		if (lbl != null) {
			return lbl.Text;
		}
		return "";
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
