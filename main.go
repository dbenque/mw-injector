package main

import (
	goflag "flag"
	"runtime"

	"github.com/dbenque/mw-injector/app"
	"github.com/spf13/pflag"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := app.NewNamespaceControllerConfig()
	config.AddFlags(pflag.CommandLine)

	// workaround for kubernetes/kubernetes/#17162
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()
	goflag.CommandLine.Parse([]string{})

	a := app.NewNamespaceController(config)
	a.Run()
}
