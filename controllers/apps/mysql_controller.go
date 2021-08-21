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

package apps

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1beta1 "github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
}

//+kubebuilder:rbac:groups=apps.shubhindia.me,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.shubhindia.me,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.shubhindia.me,resources=mysqls/finalizers,verbs=update

func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.log = log.FromContext(ctx).WithValues("Mysql", req.NamespacedName)
	r.log.Info("Started mysql reconciliation")

	funcSlice := []func(ctx context.Context, instance *v1beta1.Mysql) error{
		r.ensureDefaults, r.ensurePvc, r.ensureDeployment, r.ensureService,
	}
	instance := &v1beta1.Mysql{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		r.log.Info("Unable to fetch mysql object")
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}
	//loop through our functions slice to create respective objects
	for _, function := range funcSlice {
		err := function(ctx, instance)
		if err != nil {
			instance.Status.Status = v1beta1.MysqlStatusError
			instance.Status.Message = err.Error()
			return r.ensureStatus(ctx, instance, ctrl.Result{})
		}
	}
	//Get the deployment and check the status. If its not ready mark mysql as unready
	deployment := &v1.Deployment{}
	_ = r.Client.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	if deployment.Status.ReadyReplicas != 1 {
		instance.Status.Status = v1beta1.MySQLStatusDeploying
		instance.Status.Message = "Deployment is not ready"
		return r.ensureStatus(ctx, instance, ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second})

	}

	instance.Status.Status = v1beta1.MysqlStatusReady
	instance.Status.Message = fmt.Sprintf("Mysql instance %s is ready", instance.Name)
	return r.ensureStatus(ctx, instance, ctrl.Result{})
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1beta1.Mysql{}).
		Owns(&v1.Deployment{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}

func (r *MysqlReconciler) ensureStatus(ctx context.Context, instance *v1beta1.Mysql, result ctrl.Result) (ctrl.Result, error) {

	err := r.Status().Update(ctx, instance)
	if err != nil {
		r.log.Error(err, "Failed to update status")
		return ctrl.Result{Requeue: true}, nil
	}

	return result, nil
}
