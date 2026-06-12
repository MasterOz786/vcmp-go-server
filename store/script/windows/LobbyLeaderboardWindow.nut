// 2D leaderboard overlay (decui).

class LobbyLeaderboardController {
	component = null;
	res = null;
	visible = false;
	rootId = "safari::leaderboard";

	constructor(res) {
		this.res = res;
	}

	function hideOverlay() {
		if (this.component != null) {
			this.component.destroy();
			this.component = null;
		} else if (UI.findById(this.rootId + "::overlay") != null) {
			UI.DeleteByID(this.rootId + "::overlay");
		}
		GUI.SetMouseEnabled(false);
		this.visible = false;
	}

	function hide() {
		this.hideOverlay();
	}

	function showOverlay(escortRows, defendRows) {
		this.hideOverlay();
		this.component = LeaderboardComponent({
			id = this.rootId,
			res = this.res,
			escortRows = escortRows,
			defendRows = defendRows,
		});
		this.visible = true;
	}

	function onResize(res) {
		this.res = res;
		if (!this.visible) {
			return;
		}
		local escort = [];
		local defend = [];
		if (this.component != null) {
			escort = this.component.escortRows;
			defend = this.component.defendRows;
		}
		this.showOverlay(escort, defend);
	}
}
