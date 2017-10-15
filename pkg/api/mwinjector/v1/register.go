package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/dbenque/mw-injector/pkg/api/mwinjector"
)

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

const (
	// ResourcePlural is the id to indentify pluarals
	ResourcePlural = "mwinjectors"
	// ResourceSingular represents the id for identify singular resource
	ResourceSingular = "mwinjector"
	// ResourceKind
	ResourceKind = "MWInjector"
	// ReourceVersion
	ResourceVersion = "v1"
)

// SchemeGroupVersion is the group version used to register these objects.
var SchemeGroupVersion = schema.GroupVersion{Group: mwinjector.GroupName, Version: ResourceVersion}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&MWInjector{},
		&MWInjectorList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
