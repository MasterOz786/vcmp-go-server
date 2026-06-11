dofile("windows/RegisterWindow.nut");
dofile("windows/PacksWindow.nut");
dofile("windows/RoundScoreboardWindow.nut");

class WindowsController {
	res = null;
	registerWindow = null;
	packsWindow = null;
	roundScoreboard = null;

	constructor(res) {
		this.res = res;
		this.registerWindow = RegisterWindow(res);
		this.packsWindow = PacksWindow(res);
		this.roundScoreboard = RoundScoreboardWindow(res);
	}
}
