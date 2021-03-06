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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSAccountBindingApprovalSpec defines the desired state of AWSAccountBindingApproval
type AWSAccountBindingApprovalSpec struct {
	// +kubebuilder:validation:MaxLength=12
	// +kubebuilder:validation:Required
	// AccountID is an AWS Account ID.
	AccountID string `json:"accountID"`
	// +kubebuilder:default:false
	// Approved is whether or not to approve this binding request. To be
	// set by a cluster administrator.
	Approved bool `json:"approved"`
	// ARN is the AWS ARN to configure.
	// +kubebuilder:validation:Required
	ARN string `json:"arn"`
}

// AWSAccountBindingApprovalStatus defines the observed state of AWSAccountBindingApproval
type AWSAccountBindingApprovalStatus struct {
	// Approved indicates whether the account binding has been approved.
	Approved *bool `json:"approved"`
	// +nullable
	// ApprovedAt is the timestamp when the account binding was approved.
	ApprovedAt *string `json:"approvedAt"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Approved",type=string,JSONPath=`.status.approved`
// +kubebuilder:printcolumn:name="Approved At",type=string,JSONPath=`.status.approvedAt`

// AWSAccountBindingApproval is the Schema for the awsaccountbindingapprovals API
type AWSAccountBindingApproval struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSAccountBindingApprovalSpec   `json:"spec,omitempty"`
	Status AWSAccountBindingApprovalStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSAccountBindingApprovalList contains a list of AWSAccountBindingApproval
type AWSAccountBindingApprovalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAccountBindingApproval `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSAccountBindingApproval{}, &AWSAccountBindingApprovalList{})
}
