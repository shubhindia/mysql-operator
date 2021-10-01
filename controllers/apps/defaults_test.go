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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/shubhindia/mysql-operator/apis/apps/v1beta1"
)

var _ = Describe("MySqlDefaults", func() {

	BeforeEach(func() {
		instance = &v1beta1.Mysql{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: v1beta1.MysqlSpec{},
		}

	})
	Context("Update spec", func() {
		It("Make sure default values get added  when not provided", func() {
			ctx := context.TODO()
			reconciler.ensureDefaults(ctx, instance)
			Expect(instance.Spec.Image).To(Equal("mysql:5.6"))
			Expect(instance.Spec.PVCSpec.Name).To(Equal("mysql-pvc"))
			Expect(instance.Spec.PVCSpec.Size).To(Equal("1Gi"))
			Expect(instance.Spec.PVCSpec.StorageClassName).To(Equal("standard"))

		})
		It("Make sure that provided values are not being overwritten by default values", func() {
			instance.Name = "values-test"
			instance.Spec.Image = "mysql:7.4"
			instance.Spec.PVCSpec.Size = "2Gi"
			instance.Spec.PVCSpec.StorageClassName = "fast"
			ctx := context.TODO()
			reconciler.ensureDefaults(ctx, instance)
			Expect(instance.Spec.Image).To(Equal("mysql:7.4"))
			Expect(instance.Spec.PVCSpec.Size).To(Equal("2Gi"))
			Expect(instance.Spec.PVCSpec.StorageClassName).To(Equal("fast"))

		})

	})
})
