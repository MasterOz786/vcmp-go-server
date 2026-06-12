// Layout for store/sprites/packs.jpg (370x420, copy from CTF).
PackUI <- {
	IMG_W = 370,
	IMG_H = 420,

	function scaleFor(res) {
		local s = res.Y.tofloat() / 600.0;
		if (s < 1.0) {
			return 1.0;
		}
		if (s > 1.35) {
			return 1.35;
		}
		return s;
	},

	function frame(res) {
		local sc = scaleFor(res);
		local w = (IMG_W * sc).tointeger();
		local h = (IMG_H * sc).tointeger();
		return {
			x = (res.X / 2) - (w / 2),
			y = (res.Y / 2) - (h / 2),
			w = w,
			h = h,
			sc = sc,
		};
	},

	// Click regions over the top row of packs.jpg (first 3 columns).
	SLOTS = [
		{ nx = 0.02, ny = 0.17, nw = 0.23, nh = 0.28 },
		{ nx = 0.26, ny = 0.17, nw = 0.23, nh = 0.28 },
		{ nx = 0.50, ny = 0.17, nw = 0.23, nh = 0.28 },
	],
};
