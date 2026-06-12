dofile("decui/decui.nut");
dofile("utils/SafariConstants.nut");
dofile("utils/SafariTheme.nut");
dofile("utils/PackDefinitions.nut");
dofile("components/PackPickerComponent.nut");
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
	if (windows != null) {
		if (windows.packsWindow != null) {
			windows.packsWindow.clear();
		}
		if (windows.lobbyLeaderboard != null) {
			windows.lobbyLeaderboard.hide();
		}
		if (windows.roundScoreboard != null) {
			windows.roundScoreboard.hide();
		}
		if (windows.registerWindow != null) {
			windows.registerWindow.clear();
		}
	}
	if ("Timer" in getroottable() && Timer.Timers != null) {
		Timer.Timers.clear();
	}
	GUI.SetMouseEnabled(false);
}

function Script::ScriptLoad() {
	SafariHydraCam.init();
	SafariHydraCam.sendHello();
	packsKey <- KeyBind(0x50);

	local reg = Stream();
	reg.WriteInt(Packets.REQUEST_REGISTER_UI);
	Server.SendData(reg);
}

function Script::ScriptProcess() {
	Timer.Process();
	UI.events.scriptProcess();
}

function Server::ServerData(stream) {
	streams.process(stream);
}

function GUI::InputReturn(editbox) {
	UI.events.onInputReturn(editbox);
	if (windows.registerWindow != null && windows.registerWindow.passwordInput != null) {
		if (windows.registerWindow.passwordInput == editbox) {
			windows.registerWindow.register();
		}
	}
}

function GUI::GameResize(width, height) {
	UI.events.onGameResize();
	local v = VectorScreen(width, height);
	scoreboard.onResize(v);
	sprites.updatePositions(v);
	windows.packsWindow.updatePositions(v);
	windows.roundScoreboard.onResize(v);
	windows.lobbyLeaderboard.onResize(v);
	windows.registerWindow.updatePositions(v);
}

function GUI::ElementClick(element, mouseX, mouseY) {
	UI.events.onClick(element, mouseX, mouseY);
	clicks.handleClick(element, mouseX, mouseY);
}

function GUI::ElementFocus(element) {
	UI.events.onFocus(element);
}

function GUI::ElementBlur(element) {
	UI.events.onBlur(element);
}

function GUI::ElementHoverOver(element) {
	UI.events.onHoverOver(element);
}

function GUI::ElementHoverOut(element) {
	UI.events.onHoverOut(element);
}

function KeyBind::OnDown(bind) {
	if (SafariHydraCam.hydraKey != null && bind == SafariHydraCam.hydraKey) {
		SafariHydraCam.requestCycle();
		return;
	}
	if (packsKey != null && bind == packsKey) {
		if (windows.lobbyLeaderboard.visible) {
			windows.lobbyLeaderboard.hideOverlay();
			return;
		}
		if (windows.roundScoreboard.visible) {
			windows.roundScoreboard.hideBanner();
			return;
		}
		if (windows.packsWindow.component != null) {
			windows.packsWindow.clear();
			return;
		}
		windows.packsWindow.requestToggle();
	}
}
