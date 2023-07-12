/*
Copyright 2023 The Kubermatic Kubernetes Platform contributors.

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

package rbac

import (
	"fmt"

	"k8c.io/kubermatic/v2/pkg/crd"

	"k8s.io/apimachinery/pkg/api/meta"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// getPluralResourceName returns the plural resource name we require to generate RBAC in rbac-controller.
// If the resource type is part of our CRDs, it will be read from the compiled-in CRD definitions; otherwise
// we try to dynamically discover it (e.g. for corev1 resources).
// Relying on the compile-time CRDs allows this function to work even if the CRDs are not installed in the
// cluster, which is the case for example for Clusters, which live on seeds, but for which RBAC has to be
// generated by the master-ctrl-mgr on the master cluster.
func getPluralResourceName(restMapper meta.RESTMapper, obj ctrlruntimeclient.Object) (string, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	groups, err := crd.Groups()
	if err != nil {
		return "", fmt.Errorf("failed to get CRD groups: %w", err)
	}

	crdAvailable := false
	for _, value := range groups {
		if gvk.Group == value {
			crdAvailable = true
			break
		}
	}

	// this is static information we can discover from our CRDs.
	if crdAvailable {
		objCrd, err := crd.CRDForGVK(gvk)
		if err != nil {
			return "", fmt.Errorf("failed to get CRD for GroupVersionKind: %w", err)
		}

		return objCrd.Spec.Names.Plural, nil
	}

	// if we don't have the CRD stored, we will try to discover this dynamically.

	rmapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", fmt.Errorf("failed to get REST Mapping for '%s': %w", gvk.GroupKind().String(), err)
	}

	return rmapping.Resource.Resource, nil
}