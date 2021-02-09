package controllers

import "k8s.io/apimachinery/pkg/types"

// bindingInstanceKey returns the corresponding AWSAccountBinding's types.NamespacedName
// from the instanceKey of a given AWSAccountBindingApproval.
func bindingInstanceKey(key types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{Name: key.Name}
}

// approvalInstanceKey returns the corresponding AWSAccountBindingApproval's types.NamespacedName
// from the instanceKey of a given AWSAccountBindingRequest.
func approvalInstanceKey(key types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{Name: key.Namespace}
}
