package states

import (
	awsint "github.com/opdev/aws-account-binding-operator/api/v1alpha1"
)

// AccountBindingRequestResources stores the resources in scope for
// an account binding approval
type AccountBindingRequestResources struct {
	instance       *awsint.AWSAccountBindingRequest
	approvalStatus *awsint.AWSAccountBindingApprovalStatus
}

// NewAccountBindingRequestResources returns a new AccountBindingRequestResources struct
// which can have its state parsed.
func NewAccountBindingRequestResources(inst awsint.AWSAccountBindingRequest, appr awsint.AWSAccountBindingApprovalStatus) AccountBindingRequestResources {
	return AccountBindingRequestResources{
		instance:       &inst,
		approvalStatus: &appr,
	}
}

// Instance returns the instance of the input resource.
func (a *AccountBindingRequestResources) Instance() *awsint.AWSAccountBindingRequest {
	return a.instance
}

// ParseState returns an AccountBindingRequestState, reflecting the current state
// of the resource as determined by observing the input resource's spec.
func (a *AccountBindingRequestResources) ParseState() AccountBindingRequestState {
	return newAccountBindingRequestState(a)
}

// ARN returns the ARN associated with the request
func (a *AccountBindingRequestResources) ARN() string {
	return a.instance.Spec.ARN
}

// AccountID returns the AccountID associated with the request
func (a *AccountBindingRequestResources) AccountID() string {
	return a.instance.Spec.AccountID
}

// AccountBindingRequestState is the state of an AccountBindingApproval
// resource.
type AccountBindingRequestState struct {
	AccountBindingRequestResources
	isApproved     bool
	timeApproved   *string
	isBeingDeleted bool
	approvalKey    string
}

// newAccountBindingRequestState is the parsed state of the resources.
func newAccountBindingRequestState(resources *AccountBindingRequestResources) AccountBindingRequestState {
	var approved bool

	if resources.approvalStatus.Approved != nil && *resources.approvalStatus.Approved {
		approved = true
	}

	return AccountBindingRequestState{
		AccountBindingRequestResources: *resources,
		// isApproved:                     resources.instance.Spec.Approved,
		// timeApproved:                   resources.instance.Status.ApprovedAt,
		isBeingDeleted: !resources.instance.DeletionTimestamp.Equal(nil),
		approvalKey:    resources.instance.GetObjectMeta().GetNamespace(),
		isApproved:     approved,
	}
}

// IsApproved checks the associated approval's status to see
// if it has been approved.
func (a *AccountBindingRequestState) IsApproved() bool {
	return a.isApproved
}

// IsBeingDeleted observes the state and determines if the instance is
// being deleted by checking for a deletion timestamp.
func (a *AccountBindingRequestState) IsBeingDeleted() bool {
	return a.isBeingDeleted
}

// CurrentStatus returns an AWSAccountBindingRequestStatus based on spec.
func (a *AccountBindingRequestState) CurrentStatus() awsint.AWSAccountBindingRequestStatus {
	approved := a.IsApproved()
	return awsint.AWSAccountBindingRequestStatus{
		Approved: &approved,
	}
}
