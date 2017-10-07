package app

import (
	"github.com/spf13/pflag"
)

// Config contains configuration for workflow controller application
type Config struct {
	KubeConfigFile string
}

// NewWorkflowControllerConfig builds and returns a workflow controller Config
func NewNamespaceControllerConfig() *Config {
	return &Config{}
}

// AddFlags add cobra flags to populate Config
func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.KubeConfigFile, "kubeconfig", c.KubeConfigFile, "Location of kubecfg file for access to kubernetes master service")
}
