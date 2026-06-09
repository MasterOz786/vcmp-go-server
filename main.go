package main

import "github.com/masteroz/vcmp-go-plugin/vcmp"

func init() {
	vcmp.MetaProvider = func() vcmp.PluginMeta {
		return vcmp.PluginMeta{Name: "Safari", Version: 0x00010000}
	}
	vcmp.OnLoad = func() {
		plug = newPlugin(loadConfig())
		plug.register()
	}
}

func main() {}
