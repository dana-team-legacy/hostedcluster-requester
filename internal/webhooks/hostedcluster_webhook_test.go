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
	hyp "github.com/openshift/hypershift/api/v1alpha1"
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
		annotVal    string
		expectedVal string
		requester   string
	}{
		{
			name:        "Set annotation on HostedCluster without annotation",
			hcn:         "hc-1",
			expectedVal: "requester-1",
			requester:   "requester-1",
		},
		{
			name:        "No operation on HostedCluster with annotation",
			hcn:         "hc-2",
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
			if tc.annotVal != "" {
				hcInst.SetAnnotations(map[string]string{constants.RequesterAnnotation: tc.annotVal})
			}

			// test
			hook.handleInner(l, hcInst, tc.requester)

			// report
			g.Expect(hcInst.Annotations[constants.RequesterAnnotation]).Should(Equal(tc.expectedVal))

		})
	}
}
