// Generated file, do not modify manually!
package versioned

import (
	mwinjectorv1 "github.com/dbenque/mw-injector/pkg/client/versioned/typed/mwinjector/v1"
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	MwinjectorV1() mwinjectorv1.MwinjectorV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Mwinjector() mwinjectorv1.MwinjectorV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	mwinjectorV1 *mwinjectorv1.MwinjectorV1Client
}

// MwinjectorV1 retrieves the MwinjectorV1Client
func (c *Clientset) MwinjectorV1() mwinjectorv1.MwinjectorV1Interface {
	return c.mwinjectorV1
}

// Deprecated: Mwinjector retrieves the default version of MwinjectorClient.
// Please explicitly pick a version.
func (c *Clientset) Mwinjector() mwinjectorv1.MwinjectorV1Interface {
	return c.mwinjectorV1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.mwinjectorV1, err = mwinjectorv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.mwinjectorV1 = mwinjectorv1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.mwinjectorV1 = mwinjectorv1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
