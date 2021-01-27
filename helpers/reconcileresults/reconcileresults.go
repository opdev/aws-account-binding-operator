package reconcileresults

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DoNotRequeue returns a controller result pairing specifying not to requeue.
func DoNotRequeue() (reconcile.Result, error) { return ctrl.Result{Requeue: false}, nil }

// RequeueWithError returns a controller result pairing specifying to
// requeue with an error message.
func RequeueWithError(e error) (reconcile.Result, error) { return ctrl.Result{Requeue: true}, e }

// Requeue returns a controller result pairing specifying to
// requeue with no error message implied. This returns no error.
func Requeue() (reconcile.Result, error) { return ctrl.Result{Requeue: true}, nil }

// RequeueWithDelay returns a controller result pairing specifying to
// requeue after a delay. This returns no error.
func RequeueWithDelay(dur time.Duration) (reconcile.Result, error) {
	return ctrl.Result{Requeue: true, RequeueAfter: dur}, nil
}
