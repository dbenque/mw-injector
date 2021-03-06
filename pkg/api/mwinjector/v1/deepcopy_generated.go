// +build !ignore_autogenerated

// Generated file, do not modify manually!

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	reflect "reflect"
)

// GetGeneratedDeepCopyFuncs returns the generated funcs, since we aren't registering them.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func GetGeneratedDeepCopyFuncs() []conversion.GeneratedDeepCopyFunc {
	return []conversion.GeneratedDeepCopyFunc{
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MWInjector).DeepCopyInto(out.(*MWInjector))
			return nil
		}, InType: reflect.TypeOf(&MWInjector{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MWInjectorCondition).DeepCopyInto(out.(*MWInjectorCondition))
			return nil
		}, InType: reflect.TypeOf(&MWInjectorCondition{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MWInjectorList).DeepCopyInto(out.(*MWInjectorList))
			return nil
		}, InType: reflect.TypeOf(&MWInjectorList{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MWInjectorSpec).DeepCopyInto(out.(*MWInjectorSpec))
			return nil
		}, InType: reflect.TypeOf(&MWInjectorSpec{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*MWInjectorStatus).DeepCopyInto(out.(*MWInjectorStatus))
			return nil
		}, InType: reflect.TypeOf(&MWInjectorStatus{})},
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MWInjector) DeepCopyInto(out *MWInjector) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MWInjector.
func (in *MWInjector) DeepCopy() *MWInjector {
	if in == nil {
		return nil
	}
	out := new(MWInjector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MWInjector) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MWInjectorCondition) DeepCopyInto(out *MWInjectorCondition) {
	*out = *in
	in.LastProbeTime.DeepCopyInto(&out.LastProbeTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MWInjectorCondition.
func (in *MWInjectorCondition) DeepCopy() *MWInjectorCondition {
	if in == nil {
		return nil
	}
	out := new(MWInjectorCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MWInjectorList) DeepCopyInto(out *MWInjectorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MWInjector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MWInjectorList.
func (in *MWInjectorList) DeepCopy() *MWInjectorList {
	if in == nil {
		return nil
	}
	out := new(MWInjectorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MWInjectorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MWInjectorSpec) DeepCopyInto(out *MWInjectorSpec) {
	*out = *in
	in.MWDeployment.DeepCopyInto(&out.MWDeployment)
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		if *in == nil {
			*out = nil
		} else {
			*out = new(meta_v1.LabelSelector)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MWInjectorSpec.
func (in *MWInjectorSpec) DeepCopy() *MWInjectorSpec {
	if in == nil {
		return nil
	}
	out := new(MWInjectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MWInjectorStatus) DeepCopyInto(out *MWInjectorStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]MWInjectorCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MWInjectorStatus.
func (in *MWInjectorStatus) DeepCopy() *MWInjectorStatus {
	if in == nil {
		return nil
	}
	out := new(MWInjectorStatus)
	in.DeepCopyInto(out)
	return out
}
