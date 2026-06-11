// Place packs.jpg in store/sprites/ (copy from CTF server store/sprites/).
class SpritesController {
	packsSprite = null;
	res = null;

	constructor(res) {
		this.res = res;
	}

	function showPacksSprite() {
		if (packsSprite != null) {
			return;
		}
		local x = res.X / 2;
		local y = res.Y / 2;
		local v = VectorScreen(x - (370 / 2), y - (420 / 2));
		packsSprite = GUISprite("packs.jpg", v);
	}

	function hidePacksSprite() {
		if (packsSprite != null) {
			packsSprite = null;
		}
	}

	function updatePositions(newres) {
		this.res = newres;
		if (packsSprite == null) {
			return;
		}
		local x = newres.X / 2;
		local y = newres.Y / 2;
		packsSprite.Position = VectorScreen(x - (370 / 2), y - (420 / 2));
	}
}
