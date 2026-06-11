dofile("utils/SafariConstants.nut");
dofile("utils/PackDefinitions.nut");
dofile("utils/Timers.nut");
dofile("HydraCamera.nut");
dofile("ScoreboardHUD.nut");
dofile("SpritesController.nut");
dofile("WindowsController.nut");
dofile("StreamsController.nut");
dofile("ClickHandler.nut");

local screen = GUI.GetScreenSize();
scoreboard <- ScoreboardHUD;
scoreboard.init(screen);
sprites <- SpritesController(screen);
windows <- WindowsController(screen);
streams <- StreamsController(windows, sprites);
clicks <- ClickHandler(windows);
packsKey <- null;

function Script::ScriptUnload() {
	if (packsKey != null) {
		packsKey = null;
	}
	if (SafariHydraCam.hydraKey != null) {
		SafariHydraCam.hydraKey = null;
	}
	if (sprites != null) {
		sprites.hidePacksSprite();
	}
	if (windows != null) {
		if (windows.packsWindow != null) {
			windows.packsWindow.clear();
		}
		if (windows.roundScoreboard != null) {
			windows.roundScoreboard.hide();
		}
		if (windows.registerWindow != null) {
			windows.registerWindow.clear();
		}
	}
	Timer.Timers.clear();
	GUI.SetMouseEnabled(false);
}

function Script::ScriptLoad() {
	SafariHydraCam.init();
	SafariHydraCam.sendHello();
	packsKey <- KeyBind(0x50); // P — toggle weapon pack picker / close scoreboard

	local reg = Stream();
	reg.WriteInt(Packets.REQUEST_REGISTER_UI);
	Server.SendData(reg);
}

function Script::ScriptProcess() {
	Timer.Process();
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
	scoreboard.onResize(v);
	sprites.updatePositions(v);
	windows.packsWindow.updatePositions(v);
	windows.roundScoreboard.onResize(v);
	windows.registerWindow.updatePositions(v);
}

function GUI::ElementClick(element, mouseX, mouseY) {
	clicks.handleClick(element, mouseX, mouseY);
}

function KeyBind::OnDown(bind) {
	if (SafariHydraCam.hydraKey != null && bind == SafariHydraCam.hydraKey) {
		SafariHydraCam.requestCycle();
		return;
	}
	if (packsKey != null && bind == packsKey) {
		if (windows.roundScoreboard.visible) {
			windows.roundScoreboard.hide();
			return;
		}
		if (windows.packsWindow.canvas != null) {
			sprites.hidePacksSprite();
			windows.packsWindow.clear();
			return;
		}
		windows.packsWindow.requestToggle();
	}
}
