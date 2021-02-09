package controllers

import (
	"context"

	"github.com/opdev/aws-account-binding-operator/controllers/constants"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	reconc "github.com/opdev/aws-account-binding-operator/helpers/reconcileresults"
)

// handleFinalizer executes finalizer management
func (r *AWSAccountBindingApprovalReconciler) handleFinalizer(ctx context.Context, modifyFinalizer finalizerFunc) (*ctrl.Result, error) {
	resources, result, err := r.GetResources(ctx)
	if reconc.ShouldHaltOrRequeue(result, err) {
		return result, err
	}

	inst := resources.Instance()
	patch := client.MergeFrom(inst.DeepCopy())
	modifyFinalizer(inst, constants.Finalizer)
	if err := r.Patch(ctx, inst, patch); err != nil {
		return reconc.RequeueWithError(err)
	}

	return reconc.ContinueReconciling()
}
