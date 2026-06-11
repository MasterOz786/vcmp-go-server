class ClickHandler {
	windows = null;

	constructor(windows) {
		this.windows = windows;
	}

	function handleClick(element, mouseX, mouseY) {
		if (windows.registerWindow == null) {
			return;
		}
		if (element == windows.registerWindow.registerBtn) {
			windows.registerWindow.register();
		} else if (element == windows.registerWindow.closeBtn) {
			windows.registerWindow.closeWindow();
		}
	}
}
