class RegisterWindow {
	component = null;
	res = null;
	rootId = "safari::register";

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
		GUI.SetMouseEnabled(false);
	}

	function createWindow() {
		if (this.component != null) {
			return;
		}
		this.component = RegisterComponent({
			id = this.rootId,
			res = this.res,
		});
	}

	function updatePositions(res) {
		this.res = res;
	}
}
