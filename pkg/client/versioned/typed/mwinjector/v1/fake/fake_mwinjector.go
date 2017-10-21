// Generated file, do not modify manually!
package fake

import (
	mwinjector_v1 "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMWInjectors implements MWInjectorInterface
type FakeMWInjectors struct {
	Fake *FakeMwinjectorV1
	ns   string
}

var mwinjectorsResource = schema.GroupVersionResource{Group: "mwinjector", Version: "v1", Resource: "mwinjectors"}

var mwinjectorsKind = schema.GroupVersionKind{Group: "mwinjector", Version: "v1", Kind: "MWInjector"}

// Get takes name of the mWInjector, and returns the corresponding mWInjector object, and an error if there is any.
func (c *FakeMWInjectors) Get(name string, options v1.GetOptions) (result *mwinjector_v1.MWInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mwinjectorsResource, c.ns, name), &mwinjector_v1.MWInjector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mwinjector_v1.MWInjector), err
}

// List takes label and field selectors, and returns the list of MWInjectors that match those selectors.
func (c *FakeMWInjectors) List(opts v1.ListOptions) (result *mwinjector_v1.MWInjectorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mwinjectorsResource, mwinjectorsKind, c.ns, opts), &mwinjector_v1.MWInjectorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &mwinjector_v1.MWInjectorList{}
	for _, item := range obj.(*mwinjector_v1.MWInjectorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mWInjectors.
func (c *FakeMWInjectors) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mwinjectorsResource, c.ns, opts))

}

// Create takes the representation of a mWInjector and creates it.  Returns the server's representation of the mWInjector, and an error, if there is any.
func (c *FakeMWInjectors) Create(mWInjector *mwinjector_v1.MWInjector) (result *mwinjector_v1.MWInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mwinjectorsResource, c.ns, mWInjector), &mwinjector_v1.MWInjector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mwinjector_v1.MWInjector), err
}

// Update takes the representation of a mWInjector and updates it. Returns the server's representation of the mWInjector, and an error, if there is any.
func (c *FakeMWInjectors) Update(mWInjector *mwinjector_v1.MWInjector) (result *mwinjector_v1.MWInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mwinjectorsResource, c.ns, mWInjector), &mwinjector_v1.MWInjector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mwinjector_v1.MWInjector), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMWInjectors) UpdateStatus(mWInjector *mwinjector_v1.MWInjector) (*mwinjector_v1.MWInjector, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mwinjectorsResource, "status", c.ns, mWInjector), &mwinjector_v1.MWInjector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mwinjector_v1.MWInjector), err
}

// Delete takes name of the mWInjector and deletes it. Returns an error if one occurs.
func (c *FakeMWInjectors) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mwinjectorsResource, c.ns, name), &mwinjector_v1.MWInjector{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMWInjectors) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mwinjectorsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &mwinjector_v1.MWInjectorList{})
	return err
}

// Patch applies the patch and returns the patched mWInjector.
func (c *FakeMWInjectors) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *mwinjector_v1.MWInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mwinjectorsResource, c.ns, name, data, subresources...), &mwinjector_v1.MWInjector{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mwinjector_v1.MWInjector), err
}
