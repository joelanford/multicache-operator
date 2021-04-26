/*
Copyright 2021.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	homev1 "github.com/joelanford/multicache-operator/api/v1"
)

// CarReconciler reconciles a Car object
type CarReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=home.lanford.io,resources=cars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=home.lanford.io,resources=cars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=home.lanford.io,resources=cars/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Car object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *CarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("car", req.NamespacedName)
	r.Log.Info("reconciling")
	var car homev1.Car
	if err := r.Get(ctx, req.NamespacedName, &car); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// If we have cluster-scoped permissions, expect to see all deployments in cluster.
	// If we don't have cluster-scoped permissions, expect to see no deployments since the
	// dynamic multi-namespace cache does not yet have any added namespaces.
	depListAll := appsv1.DeploymentList{}
	if err := r.List(ctx, &depListAll); err != nil {
		return ctrl.Result{}, err
	}

	for _, dep := range depListAll.Items {
		r.Log.Info("found deployment (all)", "name", dep.Name, "namespace", dep.Namespace)
	}

	// If we have cluster-scoped permissions or permissions in namespace olm, expect to see
	// all deployments in the olm namespace.
	depListOLM := appsv1.DeploymentList{}
	if err := r.List(ctx, &depListOLM, client.InNamespace("olm")); err != nil {
		return ctrl.Result{}, err
	}

	for _, dep := range depListOLM.Items {
		r.Log.Info("found deployment (olm)", "name", dep.Name)
	}

	// If we have cluster-scoped permissions or permissions in namespace olm, expect to see
	// all deployments in the default namespace.
	depListDefault := appsv1.DeploymentList{}
	if err := r.List(ctx, &depListDefault, client.InNamespace("default")); err != nil {
		return ctrl.Result{}, err
	}

	for _, dep := range depListDefault.Items {
		r.Log.Info("found deployment (default)", "name", dep.Name)
	}

	// If we have cluster-scoped permissions, expect to see all deployments in cluster.
	// If we don't have cluster-scoped permissions, expect to see only deployments in
	// the olm and default namespaces, since previous namespace-scoped calls dynamically
	// added those namespaces to the cache.
	if err := r.List(ctx, &depListAll); err != nil {
		return ctrl.Result{}, err
	}

	for _, dep := range depListAll.Items {
		r.Log.Info("found deployment (all)", "name", dep.Name, "namespace", dep.Namespace)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&homev1.Car{}).
		Complete(r)
}
