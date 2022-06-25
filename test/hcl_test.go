package test

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/i4de/rulex/glogger"
)

type Config struct {
	IOMode  string        `hcl:"io_mode"`
	Service ServiceConfig `hcl:"service,block"`
}

type ServiceConfig struct {
	Protocol   string          `hcl:"protocol,label"`
	Type       string          `hcl:"type,label"`
	ListenAddr string          `hcl:"listen_addr"`
	Processes  []ProcessConfig `hcl:"process,block"`
}

type ProcessConfig struct {
	Type    string   `hcl:"type,label"`
	Command []string `hcl:"command"`
}

func Test_DecodeFile(t *testing.T) {
	var config Config
	err := hclsimple.DecodeFile("data/config.hcl", nil, &config)
	if err != nil {
		glogger.GLogger.Fatalf("Failed to load configuration: %s", err)
	}
	glogger.GLogger.Printf("Configuration is %#v", config)
}
