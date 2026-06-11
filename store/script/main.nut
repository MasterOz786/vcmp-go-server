dofile("utils/SafariConstants.nut");
dofile("utils/Timers.nut");
dofile("HydraCamera.nut");
dofile("WindowsController.nut");
dofile("StreamsController.nut");
dofile("ClickHandler.nut");

local screen = GUI.GetScreenSize();
windows <- WindowsController(screen);
streams <- StreamsController(windows);
clicks <- ClickHandler(windows);

function Script::ScriptLoad() {
	SafariHydraCam.init();
	SafariHydraCam.sendHello();

	local reg = Stream();
	reg.WriteInt(Packets.REQUEST_REGISTER_UI);
	Server.SendData(reg);
}

function Script::ScriptProcess() {
	Timer.Process();
	SafariHydraCam.onScriptProcess();
}

function Server::ServerData(stream) {
	streams.process(stream);
}

function GUI::InputReturn(editbox) {
	if (windows.registerWindow != null && windows.registerWindow.passwordInput != null) {
		if (windows.registerWindow.passwordInput == editbox) {
			windows.registerWindow.register();
		}
	}
}

function GUI::GameResize(width, height) {
	local v = VectorScreen(width, height);
	windows.registerWindow.updatePositions(v);
}

function GUI::ElementClick(element, mouseX, mouseY) {
	clicks.handleClick(element, mouseX, mouseY);
}

function KeyBind::OnDown(bind) {
	if (SafariHydraCam.hydraKey != null && bind == SafariHydraCam.hydraKey) {
		SafariHydraCam.cycleMode();
	}
}
