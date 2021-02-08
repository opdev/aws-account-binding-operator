package states

import (
	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
	"github.com/opdev/aws-account-binding-operator/controllers/constants"
	corev1 "k8s.io/api/core/v1"
)

// AccountBindingResources stores the resources in scope for
// a given account binding.
type AccountBindingResources struct {
	instance  *awsint.AWSAccountBinding
	namespace *corev1.Namespace
}

func NewAccountBindingResources(inst awsint.AWSAccountBinding, ns corev1.Namespace) AccountBindingResources {
	return AccountBindingResources{
		instance:  &inst,
		namespace: &ns,
	}
}

func (a *AccountBindingResources) Instance() *awsint.AWSAccountBinding { return a.instance }
func (a *AccountBindingResources) Namespace() *corev1.Namespace        { return a.namespace }
func (a *AccountBindingResources) ParseState() AccountBindingState {
	return newAccountBindingState(a.instance, a.namespace)
}

// AccountBindingState represents a given binding relationship between an instance of
// the AWSAccountBinding custom resource, a namespace, and an ACK controller configmap.
// Creation of a state should be done through the NewAccountBindingState function.
type AccountBindingState struct {
	AccountBindingResources
	requestedAccount string
	requestedARN     string
	existingAccount  string
	existingARN      string
	isBeingDeleted   bool
}

// NewAccountBindingState returns an AccountBindingState with values from the input resources.
func newAccountBindingState(inst *awsint.AWSAccountBinding, ns *corev1.Namespace) AccountBindingState {
	var existingAccount string
	if e, ok := ns.GetAnnotations()[constants.Annotation]; ok {
		existingAccount = e
	}

	return AccountBindingState{
		AccountBindingResources: AccountBindingResources{
			instance:  inst,
			namespace: ns,
		},
		requestedAccount: inst.Spec.AccountID,
		requestedARN:     inst.Spec.ARN,
		existingAccount:  existingAccount,
		existingARN:      "",
		isBeingDeleted:   !inst.DeletionTimestamp.Equal(nil),
	}
}

// getters

// RequestedAccount gets the requested aws account from the state. This is the AWS
// account to which the namespace should be bound.
func (a *AccountBindingState) RequestedAccount() string { return a.requestedAccount }

// RequestedARN gets the requested AWS ARN from the state. This is the AWS ARN
// to which the ACK controllers should assume when interacting with the requested account.
func (a *AccountBindingState) RequestedARN() string { return a.requestedARN }

// ExistingAccount gets the existing AWS account from the state. This is the existing
// annotation on a given namespace, binding the AWS account to that namespace.
func (a *AccountBindingState) ExistingAccount() string { return a.existingAccount }

// ExistingARN gets the existing AWS ARN from the state. This is the existing ARN
// stored in the ACK controller configuration at the time of the request.
func (a *AccountBindingState) ExistingARN() string { return a.existingARN }

// IsBeingDeleted observes the state and determines if the instance is
// being deleted by checking for a deletion timestamp.
func (a *AccountBindingState) IsBeingDeleted() bool {
	return a.isBeingDeleted
}

// NamespaceAnnotated observes the state and determines if the instance's
// referenced namepace is annotated with the requested AWS account.
func (a *AccountBindingState) NamespaceAnnotated() bool {
	return a.requestedAccount == a.existingAccount
}

// ConfigurationUpdated observes the state and determines if the configmap
// which houses multi-tenancy configurations for the ACK controller has been
// updated to contain the correct binding.
func (a *AccountBindingState) ConfigurationUpdated() bool {
	// placeholder
	return false
}

// CurrentStatus returns an AWSAccountBindingStatus based on the state of the
// associated namespace and the controller configmap.
func (a *AccountBindingState) CurrentStatus() awsint.AWSAccountBindingStatus {
	configState := a.ConfigurationUpdated()
	nsState := a.NamespaceAnnotated()
	return awsint.AWSAccountBindingStatus{
		ConfigurationUpdated: &configState,
		NamespaceAnnotated:   &nsState,
	}
}
