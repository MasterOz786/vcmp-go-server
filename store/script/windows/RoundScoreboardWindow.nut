// Round scoreboard controller — decui world boards + round-end overlay.

class RoundScoreboardController {
	escortBoard = null;
	defendBoard = null;
	roundEndComponent = null;
	res = null;
	visible = false;
	boardRoot = "safari::worldboard";
	roundEndRoot = "safari::roundend";

	constructor(res) {
		this.res = res;
	}

	function hideBanner() {
		if (this.roundEndComponent != null) {
			this.roundEndComponent.destroy();
			this.roundEndComponent = null;
		}
		GUI.SetMouseEnabled(false);
		this.visible = false;
	}

	function hideBoards() {
		if (this.escortBoard != null) {
			this.escortBoard.destroy();
			this.escortBoard = null;
		}
		if (this.defendBoard != null) {
			this.defendBoard.destroy();
			this.defendBoard = null;
		}
	}

	function hasLobbyBoards() {
		return this.escortBoard != null || this.defendBoard != null;
	}

	function hide() {
		this.hideBoards();
		this.hideBanner();
	}

	function playersToRows(players, team) {
		local rows = [];
		foreach (p in players) {
			if (p.team != team) {
				continue;
			}
			rows.append({
				name = p.name,
				points = p.points,
				marks = p.kills,
				wins = p.deaths,
			});
		}
		return rows;
	}

	function populateBoards(escortRows, defendRows, colsLabel = "Player          Pts   Mrk   Win") {
		this.hideBoards();
		this.escortBoard = WorldScoreboardComponent({
			id = this.boardRoot + "::escort",
			side = "escort",
			rows = escortRows,
			res = this.res,
			colsLabel = colsLabel,
		});
		this.defendBoard = WorldScoreboardComponent({
			id = this.boardRoot + "::defend",
			side = "defend",
			rows = defendRows,
			res = this.res,
			colsLabel = colsLabel,
		});
	}

	function show(winnerTeam, escortScore, defendScore, reason, players) {
		this.hide();

		local escortRows = this.playersToRows(players, Teams.ESCORT);
		local defendRows = this.playersToRows(players, Teams.DEFEND);
		this.populateBoards(escortRows, defendRows, "Player          Pts   K     D");

		this.roundEndComponent = RoundEndComponent({
			id = this.roundEndRoot,
			winnerTeam = winnerTeam,
			escortScore = escortScore,
			defendScore = defendScore,
			reason = reason,
			res = this.res,
		});
		this.visible = true;
	}

	function onResize(res) {
		this.res = res;
	}
}
