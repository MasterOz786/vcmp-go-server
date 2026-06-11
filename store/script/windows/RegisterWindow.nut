class RegisterWindow {

	static COL_SELECTED = Colour(240, 175, 238);
	static INPUT_SIZE = VectorScreen(230, 21);

	canvas = null;
	res = null;
	label = null;
	registerBtn = null;
	closeBtn = null;
	passwordInput = null;

	constructor(res) {
		this.res = res;
	}

	function clear() {
		label = null;
		canvas = null;
		registerBtn = null;
		closeBtn = null;
		passwordInput = null;
		GUI.SetMouseEnabled(false);
	}

	function closeWindow() {
		Timer.Create(this, function(text, int) {
			if (closeBtn != null) {
				GUI.SetFocusedElement(closeBtn);
			}
			clear();
		}, 1, 1, "", 0);
	}

	function updatePositions(res) {
		this.res = res;
		if (canvas != null) {
			canvas.Position = VectorScreen((res.X / 2) - (275 / 2), (res.Y / 2) - (100 / 2) + 30);
		}
	}

	function createWindow() {
		if (canvas != null) {
			return;
		}

		label = GUILabel(VectorScreen(-13, -38), COL_SELECTED, "Please enter your password: ");
		label.FontSize = 15;
		label.Size = VectorScreen(303, 230);
		label.TextAlignment = GUI_ALIGN_CENTERH;
		label.FontFlags = GUI_FFLAG_OUTLINE | GUI_FFLAG_BOLD;

		canvas = GUICanvas();
		canvas.Position = VectorScreen((res.X / 2) - (275 / 2), (res.Y / 2) - (100 / 2) + 30);
		canvas.Size = VectorScreen(275, 345);
		canvas.AddChild(label);

		passwordInput = GUIEditbox(VectorScreen(23, 103), INPUT_SIZE, Colour(75, 75, 75), "");
		passwordInput.TextColour = COL_SELECTED;
		passwordInput.AddFlags(GUI_FLAG_EDITBOX_MASKINPUT);

		registerBtn = GUIButton(VectorScreen(135, 145), VectorScreen(65, 26), COL_SELECTED, "Register");
		registerBtn.TextColour = Colour(0, 0, 0);
		registerBtn.FontName = "Verdana";
		registerBtn.FontFlags = GUI_FFLAG_BOLD;

		closeBtn = GUIButton(VectorScreen(80, 145), VectorScreen(48, 26), COL_SELECTED, "Close");
		closeBtn.TextColour = Colour(0, 0, 0);
		closeBtn.FontName = "Verdana";
		closeBtn.FontFlags = GUI_FFLAG_BOLD;

		canvas.AddChild(registerBtn);
		canvas.AddChild(closeBtn);
		canvas.AddChild(passwordInput);

		GUI.SetMouseEnabled(true);
	}

	function register() {
		if (passwordInput != null && passwordInput.Text != null && passwordInput.Text != "" && passwordInput.Text != " ") {
			local message = Stream();
			message.WriteInt(Packets.REGISTER);
			message.WriteString(passwordInput.Text);
			Server.SendData(message);
		} else {
			Console.Print("safari >> Error, empty password input!");
		}
	}
}
