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

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	appsv1 "github.com/shubhindia/mysql-operator/api/v1"
	a "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.shcn.me,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.shcn.me,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.shcn.me,resources=mysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mysql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("mysql", req.NamespacedName)
	log.Info("Started mysql reconciliation")
	mysql := &appsv1.Mysql{}
	err := r.Client.Get(ctx, req.NamespacedName, mysql)
	if err != nil {
		if errors.IsNotFound(err) { // Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("Mysql resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "No mysql instance found")
		return ctrl.Result{}, err
	}

	//Check for deployment. If it doesn't already exist, create a new one
	found := &a.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Mysql deployment not found. Creating one")
		dep := r.deployMysqlApp(mysql)
		log.Info("Creating a new Deployment ", "Deployment.Namespace ", dep.Namespace, "Deployment.Name ", dep.Name)
		err = r.Client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment ", "Deployment.Namespace ", dep.Namespace, "Deployment.Name ", dep.Name)
			return ctrl.Result{}, err
		}
		//Here I am assuming that service doesn't exist as well, so I am creating one. This needs to be handled in cleaner way
		ser := r.deployMysqlService(mysql)
		log.Info("Creating a new Service ", "Service.Namespace ", ser.Namespace, "Service.Name ", ser.Name)
		err = r.Client.Create(ctx, ser)
		if err != nil {
			log.Error("Failed to create new Service ", "Service.Namespace ", ser.Namespace, "Service.Name ", ser.Name)
			return ctrl.Result{}, err
		}
		//Here I am creating a PVC for the mysql pod to use. Each time a new pvc will be created.
		//TODO: Add a logic to create pvc only if it is not provided in the Mysql yaml.
		pvc := r.deployMysqlPVC(mysql)
		log.Info("Creating a PVC ", "PVC.Namespace ", pvc.Namespace, "PVC.Name ", pvc.Name)
		err = r.Client.Create(ctx, pvc)
		if err != nil {
			log.Error("Failed to create new PVC ", "PVC.Namespace ", pvc.Namespace, "PVC.Name ", pvc.Name)
			return ctrl.Result{}, err
		}
		// If there is no error, that means deployment was created successfully. Return and requeue
		return ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}
	//Check for desired amount of deployments
	size := mysql.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		log.Info("Changing desired size")
		err = r.Client.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment ", "Deployment.Namespace ", found.Namespace, "Deployment.Name ", found.Name)
			return ctrl.Result{}, err
		}
		//Spec updated. Return and requeue
		return ctrl.Result{Requeue: true}, nil

	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Mysql{}).
		Complete(r)
}

func (c *MysqlReconciler) deployMysqlApp(ma *appsv1.Mysql) *a.Deployment {

	replicas := ma.Spec.Size
	labels := map[string]string{"app": "mysql-containers"}
	image := ma.Spec.Image
	env := corev1.EnvVar{
		Name:  "MYSQL_ROOT_PASSWORD",
		Value: ma.Spec.Password,
	}
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ma.Name,
			Namespace: ma.Namespace,
		},
		Spec: a.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Image: image,
					Name:  ma.Name,
					Env:   []corev1.EnvVar{env},
					Ports: []corev1.ContainerPort{{
						Name:          "mysql-port",
						ContainerPort: 3306,
					}},
				}},
				},
			},
		},
	}
	ctrl.SetControllerReference(ma, dep, c.Scheme)
	return dep
}

func (c *MysqlReconciler) deployMysqlService(ma *appsv1.Mysql) *corev1.Service {
	labels := map[string]string{"app": "mysql-containers"}

	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ma.Name,
			Namespace: ma.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol: corev1.ProtocolTCP,
				Port:     3306,
				Name:     "mysql-port",
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	ctrl.SetControllerReference(ma, ser, c.Scheme)
	return ser
}

func (c *MysqlReconciler) deployMysqlPVC(ma *appsv1.Mysql) *corev1.PersistentVolumeClaim {
	/*
		TODO:
			1. Get the PVC size from Mysql app rather than hardcoding here. For now I am hardcoding it just to create initial scafold
	*/
	labels := map[string]string{"app": "mysql-containers"}
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ma.Name,
			Namespace: ma.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
		},
	}
	ctrl.SetControllerReference(ma, pvc, c.Scheme)
	return pvc
}
