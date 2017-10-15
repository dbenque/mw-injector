package v1

import (
	"k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMWInjector validates MWInjector
func ValidateMWInjector(mwi *MWInjector) field.ErrorList {
	allErrs := validation.ValidateObjectMeta(&mwi.ObjectMeta, true, validation.NameIsDNSSubdomain, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMWInjectorSpec(&(mwi.Spec), field.NewPath("spec"))...)
	return allErrs
}

// ValidateMWInjectorSpec validates MWInjectorSpec
func ValidateMWInjectorSpec(spec *MWInjectorSpec, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if spec.NamespaceSelector == nil {
		allErrs = append(allErrs, field.Required(fieldPath.Child("namespaceSelector"), ""))
	} else {
		allErrs = append(allErrs, v1validation.ValidateLabelSelector(spec.NamespaceSelector, fieldPath.Child("namespaceSelector"))...)
	}
	return allErrs
}
