/*
Copyright 2023.

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

package webhooks

import (
	"testing"

	. "github.com/onsi/gomega"
	hyp "github.com/openshift/hypershift/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/dana-team/hostedcluster-requester/internal/constants"
)

func TestAnnotateHostedClusterRequester(t *testing.T) {
	hook := &HostedClusterAnnotator{}
	l := zap.New(zap.UseDevMode(true))

	// name is the name of the test
	// hnc is the name of the HostedCluster
	// annotVal indicates if HostedCluster already has an annotation
	// expectedVal indicates what the expected value of the test is
	tests := []struct {
		name        string
		hcn         string
		hcns        string
		annotVal    string
		expectedVal string
		requester   string
	}{
		{
			name:        "Set annotation on HostedCluster without annotation",
			hcn:         "hc-1",
			hcns:        "clusters",
			expectedVal: "requester-1",
			requester:   "requester-1",
		},
		{
			name:        "No operation on HostedCluster with annotation",
			hcn:         "hc-2",
			hcns:        "clusters",
			annotVal:    "requester-x",
			expectedVal: "requester-x",
			requester:   "requester-2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			hcInst := &hyp.HostedCluster{}
			hcInst.Name = tc.hcn
			hcInst.Namespace = tc.hcns
			if tc.annotVal != "" {
				hcInst.SetAnnotations(map[string]string{constants.RequesterAnnotation: tc.annotVal})
			}

			// this is needed so that the tests don't fail on the webhook
			// for trying to access an invalid memory address or nil pointer dereference
			// TODO: Figure out a more elegant way to handle it
			hcInst.Spec.Etcd = hyp.EtcdSpec{
				ManagementType: "Managed",
				Managed: &hyp.ManagedEtcdSpec{
					Storage: hyp.ManagedEtcdStorageSpec{
						Type: "PersistentVolume",
						PersistentVolume: &hyp.PersistentVolumeEtcdStorageSpec{
							Size: &resource.Quantity{},
						},
						RestoreSnapshotURL: []string{},
					},
				},
			}

			// test
			hook.handleInner(l, hcInst, tc.requester)

			// report
			g.Expect(hcInst.Annotations[constants.RequesterAnnotation]).Should(Equal(tc.expectedVal))

		})
	}
}
