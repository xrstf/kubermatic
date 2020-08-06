/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package metricsserver

import (
	"fmt"

	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/resources/reconciling"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalNameServiceCreator returns the function to reconcile the metrics server service
func ExternalNameServiceCreator(namespace string) reconciling.NamedServiceCreatorGetter {
	return func() (string, reconciling.ServiceCreator) {
		return resources.MetricsServerExternalNameServiceName, func(se *corev1.Service) (*corev1.Service, error) {
			se.Namespace = metav1.NamespaceSystem
			se.Labels = resources.BaseAppLabels(Name, nil)

			se.Spec.Type = corev1.ServiceTypeExternalName
			se.Spec.ExternalName = fmt.Sprintf("%s.%s.svc.cluster.local", resources.MetricsServerServiceName, namespace)

			return se, nil
		}
	}
}
