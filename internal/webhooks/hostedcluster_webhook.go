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
	"context"
	"encoding/json"
	"net/http"

	"github.com/dana-team/hostedcluster-requester/internal/constants"
	"github.com/go-logr/logr"
	hyp "github.com/openshift/hypershift/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-v1beta1-hostedcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=hypershift.openshift.io,resources=hostedclusters,verbs=create;update,versions=v1beta1,name=mhostedcluster.kb.io,admissionReviewVersions=v1

// HostedClusterAnnotator adds an annotation to HostedCluster objects
type HostedClusterAnnotator struct {
	Client  client.Client
	Decoder *admission.Decoder
	Log     logr.Logger
}

// Handle adds an annotation with the name of the user who made the request
func (hook *HostedClusterAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := hook.Log.WithValues("webhook", "HostedCluster webhook", "Name", req.Name)
	log.Info("webhook request received")

	// get the incoming HostedCluster
	hostedCluster := hyp.HostedCluster{}
	if err := hook.Decoder.DecodeRaw(req.Object, &hostedCluster); err != nil {
		log.Error(err, "error decoding request")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	hook.handleInner(log, &hostedCluster, req.UserInfo.Username)
	marshaledHostedCluster, err := json.Marshal(hostedCluster)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledHostedCluster)
}

// handleInner implements the non-boilerplate logic of this mutator, allowing it to
// be more easily unit tested (ie without constructing a full admission.Request).
// Currently, we only add 'requester' annotation if the annotation is missing upon creation
func (hook *HostedClusterAnnotator) handleInner(log logr.Logger, hc *hyp.HostedCluster, requester string) {

	// add annotation if the object does not have it
	if _, hasAnnotation := hc.Annotations[constants.RequesterAnnotation]; !hasAnnotation {
		if hc.Annotations == nil {
			hc.Annotations = make(map[string]string)
		}
		log.Info("HostedCluster is missing requester annotation; adding")
		hc.Annotations[constants.RequesterAnnotation] = requester
	}

	// if spec.etcd.managed.storage.restoreSnapshotURL is nil then change it to an empty slice.
	// This is needed because in v1beta1 this field does not
	// have an omitempty tag and marshalling an empty []string slice gives nil
	// instead of empty slice and this causes an error for HostedCluster.
	// This is the only instance of the HostedCluster spec that is like this so
	// we only need to handle this here
	// TODO: Figure out a more elegant way to handle it
	if hc.Spec.Etcd.Managed.Storage.RestoreSnapshotURL == nil {
		hc.Spec.Etcd.Managed.Storage.RestoreSnapshotURL = []string{}
	}
}
