// Generated file, do not modify manually!
package v1

import (
	v1 "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1"
	scheme "github.com/dbenque/mw-injector/pkg/client/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MWInjectorsGetter has a method to return a MWInjectorInterface.
// A group's client should implement this interface.
type MWInjectorsGetter interface {
	MWInjectors(namespace string) MWInjectorInterface
}

// MWInjectorInterface has methods to work with MWInjector resources.
type MWInjectorInterface interface {
	Create(*v1.MWInjector) (*v1.MWInjector, error)
	Update(*v1.MWInjector) (*v1.MWInjector, error)
	UpdateStatus(*v1.MWInjector) (*v1.MWInjector, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.MWInjector, error)
	List(opts meta_v1.ListOptions) (*v1.MWInjectorList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.MWInjector, err error)
	MWInjectorExpansion
}

// mWInjectors implements MWInjectorInterface
type mWInjectors struct {
	client rest.Interface
	ns     string
}

// newMWInjectors returns a MWInjectors
func newMWInjectors(c *MwinjectorV1Client, namespace string) *mWInjectors {
	return &mWInjectors{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mWInjector, and returns the corresponding mWInjector object, and an error if there is any.
func (c *mWInjectors) Get(name string, options meta_v1.GetOptions) (result *v1.MWInjector, err error) {
	result = &v1.MWInjector{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mwinjectors").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MWInjectors that match those selectors.
func (c *mWInjectors) List(opts meta_v1.ListOptions) (result *v1.MWInjectorList, err error) {
	result = &v1.MWInjectorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mwinjectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mWInjectors.
func (c *mWInjectors) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mwinjectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mWInjector and creates it.  Returns the server's representation of the mWInjector, and an error, if there is any.
func (c *mWInjectors) Create(mWInjector *v1.MWInjector) (result *v1.MWInjector, err error) {
	result = &v1.MWInjector{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mwinjectors").
		Body(mWInjector).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mWInjector and updates it. Returns the server's representation of the mWInjector, and an error, if there is any.
func (c *mWInjectors) Update(mWInjector *v1.MWInjector) (result *v1.MWInjector, err error) {
	result = &v1.MWInjector{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mwinjectors").
		Name(mWInjector.Name).
		Body(mWInjector).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *mWInjectors) UpdateStatus(mWInjector *v1.MWInjector) (result *v1.MWInjector, err error) {
	result = &v1.MWInjector{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mwinjectors").
		Name(mWInjector.Name).
		SubResource("status").
		Body(mWInjector).
		Do().
		Into(result)
	return
}

// Delete takes name of the mWInjector and deletes it. Returns an error if one occurs.
func (c *mWInjectors) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mwinjectors").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mWInjectors) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mwinjectors").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mWInjector.
func (c *mWInjectors) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.MWInjector, err error) {
	result = &v1.MWInjector{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mwinjectors").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
