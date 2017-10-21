// Generated file, do not modify manually!
package v1

import (
	v1 "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1"
	"github.com/dbenque/mw-injector/pkg/client/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type MwinjectorV1Interface interface {
	RESTClient() rest.Interface
	MWInjectorsGetter
}

// MwinjectorV1Client is used to interact with features provided by the mwinjector group.
type MwinjectorV1Client struct {
	restClient rest.Interface
}

func (c *MwinjectorV1Client) MWInjectors(namespace string) MWInjectorInterface {
	return newMWInjectors(c, namespace)
}

// NewForConfig creates a new MwinjectorV1Client for the given config.
func NewForConfig(c *rest.Config) (*MwinjectorV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &MwinjectorV1Client{client}, nil
}

// NewForConfigOrDie creates a new MwinjectorV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *MwinjectorV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new MwinjectorV1Client for the given RESTClient.
func New(c rest.Interface) *MwinjectorV1Client {
	return &MwinjectorV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *MwinjectorV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
