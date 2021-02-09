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

// AWSAccountBindingRequestSpec defines the desired state of AWSAccountBindingRequest
type AWSAccountBindingRequestSpec struct {
	// +kubebuilder:validation:MaxLength=12
	// +kubebuilder:validation:Required
	// AccountID is an AWS Account ID to associate with the namespace.
	AccountID string `json:"accountID,omitempty"`
	// +kubebuilder:validation:Required
	// ARN is the AWS ARN to be assumed by ACK service controllers.
	ARN string `json:"arn,required"`
}

// AWSAccountBindingRequestStatus defines the observed state of AWSAccountBindingRequest
type AWSAccountBindingRequestStatus struct {
	// Approved indicates whether this binding request has been approved
	// by a cluster administrator.
	Approved *bool `json:"approved"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Approved",type=boolean,JSONPath=`.status.approved`

// AWSAccountBindingRequest is the Schema for the awsaccountbindingrequests API
type AWSAccountBindingRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSAccountBindingRequestSpec   `json:"spec,omitempty"`
	Status AWSAccountBindingRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSAccountBindingRequestList contains a list of AWSAccountBindingRequest
type AWSAccountBindingRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAccountBindingRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSAccountBindingRequest{}, &AWSAccountBindingRequestList{})
}
