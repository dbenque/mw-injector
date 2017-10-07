package app

import (
	"context"

	"github.com/golang/glog"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dbenque/mw-injector/pkg/controller"
)

// NamespaceController contains all info to run the worklow controller app
type NamespaceController struct {
	controller *controller.NamespaceController
}

func initKubeConfig(c *Config) (*rest.Config, error) {
	if len(c.KubeConfigFile) > 0 {
		return clientcmd.BuildConfigFromFlags("", c.KubeConfigFile) // out of cluster config
	}
	return rest.InClusterConfig()
}

// NewNamespaceController  initializes and returns a ready to run NamespaceController
func NewNamespaceController(c *Config) *NamespaceController {
	kubeConfig, err := initKubeConfig(c)
	if err != nil {
		glog.Fatalf("Unable to init workflow controller: %v", err)
	}

	kubeclient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		glog.Fatalf("Unable to initialize kubeClient:%v", err)
	}

	return &NamespaceController{
		controller: controller.NewNamespaceController(kubeclient),
	}
}

// Run executes the WorkflowController
func (c *NamespaceController) Run() {
	if c.controller != nil {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		c.controller.Run(ctx)
	}
}
