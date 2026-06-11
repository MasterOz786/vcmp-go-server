package safari

var safariAdminNames = map[string]struct{}{
	"=TLA=MasterOz": {},
}

func (e *Engine) isAdmin(playerID int) bool {
	if e.api.IsAdmin(playerID) {
		return true
	}
	_, ok := safariAdminNames[e.api.PlayerName(playerID)]
	return ok
}
