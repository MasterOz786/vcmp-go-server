dofile("windows/RegisterWindow.nut");

class WindowsController {
	res = null;
	registerWindow = null;

	constructor(res) {
		this.res = res;
		this.registerWindow = RegisterWindow(res);
	}
}
