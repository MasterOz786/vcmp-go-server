class ClickHandler {
	windows = null;

	constructor(windows) {
		this.windows = windows;
	}

	function handleClick(element, mouseX, mouseY) {
		if (windows.registerWindow != null && windows.registerWindow.registerBtn != null) {
			if (element == windows.registerWindow.registerBtn) {
				windows.registerWindow.register();
				return;
			}
			if (element == windows.registerWindow.closeBtn) {
				windows.registerWindow.closeWindow();
				return;
			}
		}

		if (windows.packsWindow != null && windows.packsWindow.packButtons != null) {
			for (local idx = 0; idx < windows.packsWindow.packButtons.len(); idx++) {
				if (element == windows.packsWindow.packButtons[idx]) {
					windows.packsWindow.selectPack(idx + 1);
					return;
				}
			}
			if (element.Text != null && element.Text.find("pack") != null) {
				local t = element.Text;
				local n = t.slice(t.find("pack") + 4, t.find("pack") + 5).tointeger();
				if (n >= 1 && n <= 3) {
					windows.packsWindow.selectPack(n);
					return;
				}
			}
		}
	}
}
