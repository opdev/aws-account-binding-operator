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

// AWSAccountBindingSpec defines the desired state of AWSAccountBinding
type AWSAccountBindingSpec struct {
	// +kubebuilder:validation:required
	// AccountID is the AWS AccountID to Bind
	AccountID string `json:"accountID"`
	// +kubebuilder:validation:required
	// ARN is the AWS ARN to configure.
	ARN string `json:"arn"`
}

// AWSAccountBindingStatus defines the observed state of AWSAccountBinding
type AWSAccountBindingStatus struct {
	// NamespaceAnnotated identifies whether the Namespaced has been labeled
	NamespaceAnnotated *bool `json:"namespaceAnnotated"`
	// ConfigurationUpdated identifies whether the ACK configmap has been updated
	ConfigurationUpdated *bool `json:"configurationUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Namespaced Annotated",type=boolean,JSONPath=`.status.namespaceAnnotated`
// +kubebuilder:printcolumn:name="Configuration Updated",type=boolean,JSONPath=`.status.configurationUpdated`

// AWSAccountBinding is the Schema for the awsaccountbindings API
type AWSAccountBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSAccountBindingSpec   `json:"spec,omitempty"`
	Status AWSAccountBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSAccountBindingList contains a list of AWSAccountBinding
type AWSAccountBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAccountBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSAccountBinding{}, &AWSAccountBindingList{})
}
