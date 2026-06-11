dofile("windows/RegisterWindow.nut");
dofile("windows/PacksWindow.nut");
dofile("windows/RoundScoreboardWindow.nut");
dofile("windows/LobbyLeaderboardWindow.nut");

class WindowsController {
	res = null;
	registerWindow = null;
	packsWindow = null;
	roundScoreboard = null;
	lobbyLeaderboard = null;

	constructor(res) {
		this.res = res;
		this.registerWindow = RegisterWindow(res);
		this.packsWindow = PacksWindow(res);
		this.roundScoreboard = RoundScoreboardController(res);
		this.lobbyLeaderboard = LobbyLeaderboardController(res);
	}
}
