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
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (r *MysqlReconciler) ensureDeployment(ctx context.Context, instance *v1beta1.Mysql) error {

	//Store the secret in k8s secret so can be used later
	secretName := instance.Name + "-user-password"
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: instance.Namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"password": []byte(StringWithCharset(10, charset)),
		},
	}
	err := r.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: instance.Namespace}, secret)
	if err != nil {

		if k8serrors.IsNotFound(err) {
			//creating pvc
			err = ctrl.SetControllerReference(instance, secret, r.Scheme)
			if err != nil {
				return errors.Wrapf(err, "Error setting owner reference")
			}
			err = r.Client.Create(ctx, secret)
			if err != nil {
				return errors.Wrapf(err, "Error creating a secret")
			}

			return nil
		}
		return errors.Wrapf(err, "Error getting secret")
	}
	//Create mysql deployment here
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},

		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "mysql"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mysql"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: instance.Spec.Image,
						Name:  "mysql",

						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: string(secret.Data["password"]),
							},
						},

						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "mysql-persistent-storage",
								MountPath: "/var/lib/mysql",
							},
						},
					},
					},

					Volumes: []corev1.Volume{

						{
							Name: "mysql-persistent-storage",
							VolumeSource: corev1.VolumeSource{

								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: instance.Spec.PVCSpec.Name,
								},
							},
						},
					},
				},
			},
		},
	}
	err = r.Client.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: instance.Namespace}, deployment)
	if err != nil {

		if k8serrors.IsNotFound(err) {
			//creating pvc
			err = ctrl.SetControllerReference(instance, deployment, r.Scheme)
			if err != nil {
				return errors.Wrapf(err, "Error setting owner reference")
			}
			err = r.Client.Create(ctx, deployment)
			if err != nil {
				return errors.Wrapf(err, "Error creating a deployment")
			}

			return nil
		}
		return errors.Wrapf(err, "Error getting deployment")
	}

	return nil
}
