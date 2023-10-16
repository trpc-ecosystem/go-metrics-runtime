//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

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
