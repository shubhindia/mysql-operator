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

	"github.com/pkg/errors"
	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *MysqlReconciler) ensurePvc(ctx context.Context, instance *v1beta1.Mysql) error {

	mysqlPVC := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.PVCSpec.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse(instance.Spec.PVCSpec.Size),
				},
			},
			StorageClassName: &instance.Spec.PVCSpec.StorageClassName,
		},
	}

	err := r.Client.Get(ctx, types.NamespacedName{Name: mysqlPVC.Name, Namespace: instance.Namespace}, mysqlPVC)
	if err != nil {

		if k8serrors.IsNotFound(err) {
			//creating pvc
			err = ctrl.SetControllerReference(instance, mysqlPVC, r.Scheme)
			if err != nil {
				return errors.Wrapf(err, "Error setting owner reference")
			}
			err = r.Client.Create(ctx, mysqlPVC)
			if err != nil {
				return errors.Wrapf(err, "Error creating a secret")
			}

			return nil
		}
		return errors.Wrapf(err, "Error getting pvc")
	}

	return nil
}
