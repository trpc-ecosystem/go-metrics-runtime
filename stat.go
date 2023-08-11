package runtime

import (
	"trpc.group/trpc-go/trpc-go/plugin"
)

// Plugin
const (
	pluginType = "runtime"
	pluginName = "stat"
)

var extraConf *Config

func init() {
	// register plugin
	plugin.Register(pluginName, &Plugin{})
}

// Config is a struct of configurations
type Config struct {
	Disable bool `json:"disable"` // disable plugin
}

// Plugin is a stat plugin
type Plugin struct{}

// Type returns the plugin type
func (m *Plugin) Type() string {
	return pluginType
}

// Setup initialize this plugin
func (m *Plugin) Setup(name string, configDesc plugin.Decoder) (err error) {
	var config Config
	if err = configDesc.Decode(&config); err != nil {
		return
	}
	extraConf = &config
	return nil
}

// GetExtraConf returns *Config
func GetExtraConf() *Config {
	if extraConf == nil {
		return &Config{}
	}
	return extraConf
}
