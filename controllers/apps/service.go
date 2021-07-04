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
	ctrl "sigs.k8s.io/controller-runtime"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *MysqlReconciler) ensureService(ctx context.Context, instance *v1beta1.Mysql) error {

	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "mysql"},

			Ports: []corev1.ServicePort{
				{
					Port: 3306,
					Name: "mysql",
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	if err := ctrl.SetControllerReference(instance, service, r.Scheme); err != nil {
		return errors.Wrapf(err, "Error setting owner reference")
	}
	if err := r.Client.Create(ctx, service); err != nil {
		return errors.Wrapf(err, "Error creating a service")
	}
	return nil
}
