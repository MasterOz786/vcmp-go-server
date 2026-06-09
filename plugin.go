package main

type Plugin struct {
	demo *Demo
}

func newPlugin(cfg Config) *Plugin {
	d := newDemo(cfg)
	d.register()
	return &Plugin{demo: d}
}
