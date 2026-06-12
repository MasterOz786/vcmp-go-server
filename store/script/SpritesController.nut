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
		local f = PackUI.frame(res);
		packsSprite = GUISprite("packs.jpg", VectorScreen(f.x, f.y));
		packsSprite.Size = VectorScreen(f.w, f.h);
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
		local f = PackUI.frame(newres);
		packsSprite.Position = VectorScreen(f.x, f.y);
		packsSprite.Size = VectorScreen(f.w, f.h);
	}
}
