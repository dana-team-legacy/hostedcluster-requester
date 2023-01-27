# hostedcluster-requester
This repository includes a low-level API implementation of a Mutating Webhook that adds an annotation to a `HostedCluster` object. The annotation that is added is `rks.dana.io/requester: <REQUESTER>` when `<REQUESTER>` is the name of the user that made the request.

## RKS
This webhook is part of the Run Kubernetes Service (Kubernetes as a Service solution) eco-system and is needed to grant permissions to the relevant user on its `Hypershift` cluster.

## Getting Started
To run this webhook, you need an `OpenShift` cluster with `Hypershift` installed.

### Running on the cluster
1. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/hostedcluster-requester:tag
```
	
2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/hostedcluster-requester:tag
```

## License

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

