class StreamsController {
	windows = null;

	constructor(windows) {
		this.windows = windows;
	}

	function process(stream) {
		local type = stream.ReadInt();

		if (type == Packets.HIDE_REGISTER) {
			if (windows.registerWindow != null && windows.registerWindow.registerBtn != null) {
				GUI.SetFocusedElement(windows.registerWindow.registerBtn);
				windows.registerWindow.clear();
			}
		} else if (type == Packets.SHOW_REGISTER) {
			windows.registerWindow.createWindow();
		} else if (type == Packets.HYDRA_CAM) {
			SafariHydraCam.applyMode(stream.ReadInt(), stream.ReadInt());
		}
	}
}
