/*
Copyright 2020 The Crossplane Authors.

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

package v1alpha1

import (
	"context"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/reference"

	clusterv1alpha1 "github.com/crossplane-contrib/provider-argocd/apis/cluster/v1alpha1"
	repositoriesv1alpha1 "github.com/crossplane-contrib/provider-argocd/apis/repositories/v1alpha1"
)

// ResolveReferences of this Cluster
func (mg *Project) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	// Resolve spec.forProvider.roleArn
	rsp, err := r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: mg.Spec.ForProvider.SourceRepos,
		References:    mg.Spec.ForProvider.SourceReposRefs,
		Selector:      mg.Spec.ForProvider.SourceReposSelector,
		To:            reference.To{Managed: &repositoriesv1alpha1.Repository{}, List: &repositoriesv1alpha1.RepositoryList{}},
		Extract:       reference.ExternalName(),
	})
	if err != nil {
		return errors.Wrap(err, "spec.forProvider.SourceRepos")
	}
	mg.Spec.ForProvider.SourceRepos = rsp.ResolvedValues
	mg.Spec.ForProvider.SourceReposRefs = rsp.ResolvedReferences

	for i := range mg.Spec.ForProvider.Destinations {
		rsp, err := r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: ptrToString(mg.Spec.ForProvider.Destinations[i].Server),
			Reference:    mg.Spec.ForProvider.Destinations[i].ServerRef,
			Selector:     mg.Spec.ForProvider.Destinations[i].ServerSelector,
			To:           reference.To{Managed: &clusterv1alpha1.Cluster{}, List: &clusterv1alpha1.ClusterList{}},
			Extract:      reference.ExternalName(),
		})
		if err != nil {
			return errors.Wrap(err, "spec.forProvider.Destinations")
		}
		mg.Spec.ForProvider.Destinations[i].Server = &rsp.ResolvedValue
		mg.Spec.ForProvider.Destinations[i].ServerRef = rsp.ResolvedReference
	}

	return nil
}

func ptrToString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
