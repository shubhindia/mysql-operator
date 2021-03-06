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

	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *MysqlReconciler) ensureDefaults(ctx context.Context, instance *v1beta1.Mysql) (ctrl.Result, error) {

	if instance.Spec.Image == "" {
		instance.Spec.Image = "mysql:5.6"
	}
	if instance.Spec.PVCSpec.Name == "" {
		instance.Spec.PVCSpec.Name = "mysql-pvc"
	}
	if instance.Spec.PVCSpec.Size == "" {
		instance.Spec.PVCSpec.Size = "1Gi"
	}
	if instance.Spec.PVCSpec.StorageClassName == "" {
		instance.Spec.PVCSpec.StorageClassName = "standard"
	}

	return ctrl.Result{}, nil
}
