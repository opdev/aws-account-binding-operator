// +build !ignore_autogenerated

/*
Copyright 2021 The OpDev Team.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBinding) DeepCopyInto(out *AWSAccountBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBinding.
func (in *AWSAccountBinding) DeepCopy() *AWSAccountBinding {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingApproval) DeepCopyInto(out *AWSAccountBindingApproval) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingApproval.
func (in *AWSAccountBindingApproval) DeepCopy() *AWSAccountBindingApproval {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingApproval)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBindingApproval) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingApprovalList) DeepCopyInto(out *AWSAccountBindingApprovalList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AWSAccountBindingApproval, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingApprovalList.
func (in *AWSAccountBindingApprovalList) DeepCopy() *AWSAccountBindingApprovalList {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingApprovalList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBindingApprovalList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingApprovalSpec) DeepCopyInto(out *AWSAccountBindingApprovalSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingApprovalSpec.
func (in *AWSAccountBindingApprovalSpec) DeepCopy() *AWSAccountBindingApprovalSpec {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingApprovalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingApprovalStatus) DeepCopyInto(out *AWSAccountBindingApprovalStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingApprovalStatus.
func (in *AWSAccountBindingApprovalStatus) DeepCopy() *AWSAccountBindingApprovalStatus {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingApprovalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingList) DeepCopyInto(out *AWSAccountBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AWSAccountBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingList.
func (in *AWSAccountBindingList) DeepCopy() *AWSAccountBindingList {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingRequest) DeepCopyInto(out *AWSAccountBindingRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingRequest.
func (in *AWSAccountBindingRequest) DeepCopy() *AWSAccountBindingRequest {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBindingRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingRequestList) DeepCopyInto(out *AWSAccountBindingRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AWSAccountBindingRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingRequestList.
func (in *AWSAccountBindingRequestList) DeepCopy() *AWSAccountBindingRequestList {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AWSAccountBindingRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingRequestSpec) DeepCopyInto(out *AWSAccountBindingRequestSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingRequestSpec.
func (in *AWSAccountBindingRequestSpec) DeepCopy() *AWSAccountBindingRequestSpec {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingRequestStatus) DeepCopyInto(out *AWSAccountBindingRequestStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingRequestStatus.
func (in *AWSAccountBindingRequestStatus) DeepCopy() *AWSAccountBindingRequestStatus {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingSpec) DeepCopyInto(out *AWSAccountBindingSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingSpec.
func (in *AWSAccountBindingSpec) DeepCopy() *AWSAccountBindingSpec {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAccountBindingStatus) DeepCopyInto(out *AWSAccountBindingStatus) {
	*out = *in
	if in.NamespaceAnnotated != nil {
		in, out := &in.NamespaceAnnotated, &out.NamespaceAnnotated
		*out = new(bool)
		**out = **in
	}
	if in.ConfigurationUpdated != nil {
		in, out := &in.ConfigurationUpdated, &out.ConfigurationUpdated
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAccountBindingStatus.
func (in *AWSAccountBindingStatus) DeepCopy() *AWSAccountBindingStatus {
	if in == nil {
		return nil
	}
	out := new(AWSAccountBindingStatus)
	in.DeepCopyInto(out)
	return out
}
