package controllers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

// subreconcilerFuncs are functions that are called by Reconcile() functions
// in an ordered fashion. Returning a ctrl.Result with a value of nil
// indicates that the Reconcile() function should continue reconciling.
// Any other returned ctrl.Result indicates to the Reconcile() function
// that reconciliation should halt.
type subreconcilerFuncs = func(context.Context) (*ctrl.Result, error)
