// Package v1 is the version 1 of the MWInjector ApI
// +k8s:deepcopy-gen=package
package v1

import "reflect"

// IsMWInjectorDefaulted check wether
func IsMWInjectorDefaulted(w *MWInjector) bool {
	defaultedMWInjector := DefaultMWInjector(w)
	return reflect.DeepEqual(w.Spec, defaultedMWInjector.Spec)
}

// DefaultMWInjector defaults MWInjector
func DefaultMWInjector(undefaultedMWInjector *MWInjector) *MWInjector {
	w := undefaultedMWInjector.DeepCopy()
	return w
}
