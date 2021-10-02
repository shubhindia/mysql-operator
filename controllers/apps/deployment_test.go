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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MySqlDefaults", func() {

	BeforeEach(func() {
		instance = &v1beta1.Mysql{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deployment-test",
				Namespace: "default",
				UID:       types.UID("deployment-test"),
			},
			Spec: v1beta1.MysqlSpec{},
		}

	})
	Context("Update spec", func() {
		It("Make sure deployment is created", func() {
			ctx := context.TODO()
			res, err := reconciler.ensureDeployment(ctx, instance)
			Expect(res).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).To(BeNil())
			Expect(res).To(Equal(ctrl.Result{Requeue: true}))
			Expect(err).To(BeNil())

			By("Checking id deployment is created")
			deployment := &appsv1.Deployment{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      instance.Name,
				Namespace: instance.Namespace,
			}, deployment)).To(Succeed())

		})

	})
})
