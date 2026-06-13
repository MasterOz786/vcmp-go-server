package apidef

import "fmt"

func PackCommandHelpLine() string {
	return fmt.Sprintf("/pack 1-%d - choose loadout (P or keys 1-%d)", MaxPack, MaxPack)
}

func SkinCommandHelpLine() string {
	return fmt.Sprintf("/skin 1-%d - choose character skin (spawn screen class)", MaxSkin)
}
