// Generated file, do not modify manually!
package fake

import (
	v1 "github.com/dbenque/mw-injector/pkg/client/versioned/typed/mwinjector/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMwinjectorV1 struct {
	*testing.Fake
}

func (c *FakeMwinjectorV1) MWInjectors(namespace string) v1.MWInjectorInterface {
	return &FakeMWInjectors{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMwinjectorV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
