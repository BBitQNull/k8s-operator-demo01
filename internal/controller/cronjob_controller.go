/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/BBitQNull/CronJob/api/v1"
)

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.bbitqnull.com,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.bbitqnull.com,resources=cronjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.bbitqnull.com,resources=cronjobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CronJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// TODO(user): your logic here
	var counter examplev1.CronJob
	if err := r.Get(ctx, req.NamespacedName, &counter); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling Counter",
		"name", counter.Name,
		"namespace", counter.Namespace,
		"target", counter.Spec.Target,
		"current", counter.Status.Current)

	if counter.Status.Current < counter.Spec.Target {
		counter.Status.Current++

		// 更新 status （推荐用 Status().Update 而不是 Update）
		if err := r.Status().Update(ctx, &counter); err != nil {
			logger.Error(err, "unable to update Counter status")
			return ctrl.Result{}, err
		}

		// 1 秒后再调和一次，直到达到目标
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}
	logger.Info("Counter reached target", "current", counter.Status.Current)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.CronJob{}).
		Named("cronjob").
		Complete(r)
}
