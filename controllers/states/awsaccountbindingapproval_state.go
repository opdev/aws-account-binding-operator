package states

import (
	"time"

	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
)

// AccountBindingApprovalResources stores the resources in scope for
// an account binding approval
type AccountBindingApprovalResources struct {
	instance *awsint.AWSAccountBindingApproval
}

// NewAccountBindingApprovalResources returns a new AccountBindingApprovalResources struct
// which can have its state parsed.
func NewAccountBindingApprovalResources(inst awsint.AWSAccountBindingApproval) AccountBindingApprovalResources {
	return AccountBindingApprovalResources{
		instance: &inst,
	}
}

// Instance returns the instance of the input resource.
func (a *AccountBindingApprovalResources) Instance() *awsint.AWSAccountBindingApproval {
	return a.instance
}

// ParseState returns an AccountBindingApprovalState, reflecting the current state
// of the resource as determined by observing the input resource's spec.
func (a *AccountBindingApprovalResources) ParseState() AccountBindingApprovalState {
	return newAccountBindingApprovalState(a)
}

// ARN returns the ARN associated with the request
func (a *AccountBindingApprovalResources) ARN() string {
	return a.instance.Spec.ARN
}

// AccountID returns the AccountID associated with the request
func (a *AccountBindingApprovalResources) AccountID() string {
	return a.instance.Spec.AccountID
}

// AccountBindingApprovalState is the state of an AccountBindingApproval
// resource.
type AccountBindingApprovalState struct {
	AccountBindingApprovalResources
	isApproved     bool
	timeApproved   *string
	isBeingDeleted bool
	bindingKey     string
}

// newAccountBindingApprovalState is the parsed state of the resources.
func newAccountBindingApprovalState(resources *AccountBindingApprovalResources) AccountBindingApprovalState {
	return AccountBindingApprovalState{
		AccountBindingApprovalResources: *resources,
		isApproved:                      resources.instance.Spec.Approved,
		timeApproved:                    resources.instance.Status.ApprovedAt,
		isBeingDeleted:                  !resources.instance.DeletionTimestamp.Equal(nil),
		bindingKey:                      resources.instance.GetObjectMeta().GetNamespace(),
	}
}

// IsApproved checks the resource's spec and determines if it has
// been approved.
func (a *AccountBindingApprovalState) IsApproved() bool {
	return a.isApproved
}

// HasApprovalTimestamp returns true if the approval timestamp is not nil,
// indicating it has been set.
func (a *AccountBindingApprovalState) HasApprovalTimestamp() bool {
	return a.timeApproved != nil
}

// IsBeingDeleted observes the state and determines if the instance is
// being deleted by checking for a deletion timestamp.
func (a *AccountBindingApprovalState) IsBeingDeleted() bool {
	return a.isBeingDeleted
}

// CurrentStatus returns an AWSAccountBindingApprovalStatus based on spec.
func (a *AccountBindingApprovalState) CurrentStatus() awsint.AWSAccountBindingApprovalStatus {
	var approvedTimestamp *string
	var approved = a.IsApproved()

	if approved {
		// if the current state is approved, then we try to assume
		// that the timestamp is already set as well.
		approvedTimestamp = a.instance.Status.ApprovedAt
	}

	if !a.HasApprovalTimestamp() && approved {
		// If we don't have an approval timestamp in existing status,
		// and the resource has been approved, we'll set one

		// TODO this could cause issues if a status block is missing
		// and a subsequent reconciliation changes the timestamp as a result.
		// It would alter the "actual" approval timestamp.
		// Would it make sense to capture the timestmap elsewhere and
		// use that as the source of truth?

		ts := time.Now().UTC().String()
		approvedTimestamp = &ts
	}

	return awsint.AWSAccountBindingApprovalStatus{
		Approved:   &approved,
		ApprovedAt: approvedTimestamp,
	}
}
