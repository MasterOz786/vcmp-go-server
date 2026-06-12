// Account registration dialog (decui).

class RegisterComponent extends Component {
	className = "RegisterComponent";

	res = null;
	rootId = null;

	constructor(o) {
		this.id = o.id;
		this.rootId = o.id + "::overlay";
		this.res = o.res;
		base.constructor(this.id, o);
		this.metadata.list = "canvas";
		this.build();
	}

	function submitRegister() {
		local input = UI.Editbox(id + "::password");
		if (input == null || input.Text == null || input.Text == "" || input.Text == " ") {
			Console.Print("safari >> Error, empty password input!");
			return;
		}
		local message = Stream();
		message.WriteInt(Packets.REGISTER);
		message.WriteString(input.Text);
		Server.SendData(message);
	}

	function build() {
		local self = this;
		local panel = UI.Canvas({
			id = id,
			align = "center",
			Size = VectorScreen(320, 200),
			Colour = SafariTheme.PANEL,
			children = [
				UI.Label({
					id = id + "::title",
					align = "center",
					Position = VectorScreen(0, 16),
					Size = VectorScreen(320, 24),
					Text = "Register your account",
					TextColour = SafariTheme.SELECTED,
					FontSize = 16,
					FontFlags = GUI_FFLAG_BOLD | GUI_FFLAG_OUTLINE,
				}),
				UI.Label({
					id = id + "::prompt",
					align = "center",
					Position = VectorScreen(0, 44),
					Size = VectorScreen(320, 18),
					Text = "Enter a password:",
					TextColour = SafariTheme.MUTED,
					FontSize = 12,
					FontFlags = GUI_FFLAG_OUTLINE,
				}),
				UI.Editbox({
					id = id + "::password",
					Position = VectorScreen(40, 72),
					Size = VectorScreen(240, 24),
					Colour = Colour(40, 40, 40),
					TextColour = SafariTheme.SELECTED,
					flags = GUI_FLAG_EDITBOX_MASKINPUT,
					onInputReturn = function() { self.submitRegister(); },
				}),
				UI.Button({
					id = id + "::close",
					Text = "Close",
					TextColour = SafariTheme.TEXT,
					Position = VectorScreen(72, 120),
					Size = VectorScreen(72, 28),
					Colour = SafariTheme.BUTTON,
					onClick = function() { self.destroy(); GUI.SetMouseEnabled(false); },
				}),
				UI.Button({
					id = id + "::register",
					Text = "Register",
					TextColour = SafariTheme.TEXT,
					Position = VectorScreen(168, 120),
					Size = VectorScreen(88, 28),
					Colour = SafariTheme.SELECTED,
					onClick = function() { self.submitRegister(); },
				}),
			],
		});
		panel.addBorders({ size = 2, color = SafariTheme.BORDER });

		UI.Canvas({
			id = rootId,
			align = "center",
			Size = VectorScreen(res.X, res.Y),
			Position = VectorScreen(0, 0),
			Colour = SafariTheme.OVERLAY,
			children = [panel],
		});

		GUI.SetMouseEnabled(true);
	}

	function destroy() {
		if (UI.findById(rootId) != null) {
			UI.DeleteByID(rootId);
		}
	}
}
